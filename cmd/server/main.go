package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"text/template"

	"go_app/pkg/cache"
	"go_app/pkg/database"
	"go_app/pkg/database/postgres"
	"go_app/pkg/order"
	"go_app/pkg/stansub"

	"github.com/joho/godotenv"
)

var (
	db         *sql.DB
	newCache   *cache.Cache
	signalChan = make(chan os.Signal, 1)
)

func listenSub(wg *sync.WaitGroup) {
	defer wg.Done()
	for orderSub := range stansub.JsonData {
		log.Printf("Received new order: ID = %s\n", orderSub.OrderID)
		go func(orderSub *order.Order) {
			newCache.Set(orderSub.OrderID, orderSub, 0)
			database.InsertOrderAndAll(orderSub, db)
		}(orderSub)
	}
}

// First ctrl+c stopping stan subscription, second command ctrl+c stop server

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// env variables
	port := os.Getenv("SERVER_PORT")
	host := os.Getenv("SERVER_HOST")
	serverUrl := host + ":" + port
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("SSLMODE")
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, dbname, sslmode)
	db = database.DbConnect(connStr)
	defaultExpiration, cleanupInterval := cache.ParseCleanupExpiration()
	newCache = cache.New(defaultExpiration, cleanupInterval)

	orders, err := postgres.SelectOrders(db, "", "", 20, 0)
	if err != nil {
		log.Fatal(err)
	}
	newCache.SetAllOrders(orders, 0)
	wg := new(sync.WaitGroup)
	sc, sub := stansub.SubStart()
	wg.Add(1)
	go listenSub(wg)
	wg.Add(1)
	// Ctrl + C gracious exit: stan close and unsubscribe
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for range signalChan {
			close(signalChan)
			signal.Stop(signalChan)
			stansub.SubClose(sc, sub)
		}
	}(wg)
	signal.Notify(signalChan, os.Interrupt)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		http.HandleFunc("/", handlerOrders)
		http.HandleFunc("/order_id", handlerOrderId)
		err = http.ListenAndServe(serverUrl, nil)
		if err != nil {
			fmt.Printf("Server start up failure %v", err)
			return
		}
	}(wg)
	wg.Wait()
}

func handlerOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		orders, err := getOrders()
		if err != nil {
			log.Println(err)
			w.Write([]byte("Orders not found"))
			return
		}
		templ, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Fatal(err)
		}
		err = templ.Execute(w, orders)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Responded to request: orders")
	}
}

func handlerOrderId(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet || r.Method == http.MethodPost {
		val := r.FormValue("order")
		order, err := getOrder(val)
		if err != nil {
			resp := fmt.Sprintf("Order id wrong or not found  %s", val)
			log.Println(resp)
			w.Write([]byte(resp))
			return
		}
		templ, err := template.ParseFiles("templates/order.html")
		if err != nil {
			log.Fatal(err)
		}
		err = templ.Execute(w, order)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Responded to request: order_id %s\n", val)
	}
}

func getOrder(orderID string) (*order.Order, error) {
	order, ok := newCache.Get(orderID)
	if !ok {
		orderDB, err := postgres.SelectOrder(orderID, db)
		if err != nil {
			return nil, err
		}
		newCache.Set(orderDB.OrderID, orderDB, 0)
		return orderDB, nil
	}
	return order, nil
}

func getOrders() ([]*order.Order, error) {
	orders, ok := newCache.GetAllOrders()
	if !ok {
		ordersDB, err := postgres.SelectOrders(db, "", "", 0, 0)
		if err != nil {
			return nil, err
		}
		newCache.SetAllOrders(ordersDB, 0)
		return ordersDB, nil
	}
	return orders, nil
}
