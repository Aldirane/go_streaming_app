package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

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
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("SSLMODE")
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, dbname, sslmode)
	db = database.DbConnect(connStr)
	defaultExpiration, cleanupInterval := cache.ParseCleanupExpiration()
	newCache = cache.New(defaultExpiration, cleanupInterval)

	orders, err := postgres.SelectOrders(db)
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
		http.HandleFunc("/orders", handlerOrders)
		http.HandleFunc("/order_id", handlerOrderId)
		err = http.ListenAndServe(":5555", nil)
		if err != nil {
			fmt.Printf("Server start up failure %v", err)
			return
		}
	}(wg)
	wg.Wait()
}

func handlerOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		orders, ok := newCache.GetAllOrders()
		if !ok {
			log.Println("Orders not found")
			w.Write([]byte("Orders not found"))
			return
		}
		jsonData, err := json.Marshal(orders)
		if err != nil {
			log.Println(err)
			w.Write([]byte("Server couldn't send jsonData"))
			return
		}
		w.Write([]byte(jsonData))
		log.Println("Responded to request: orders")
	}
}

func handlerOrderId(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet || r.Method == http.MethodPost {
		val := r.FormValue("order_id")
		order, ok := newCache.Get(val)
		if !ok {
			resp := fmt.Sprintf("Order id wrong or not found  %s", val)
			log.Println(resp)
			w.Write([]byte(resp))
			return
		}
		jsonData, err := json.Marshal(order)
		if err != nil {
			log.Println(err)
			w.Write([]byte("Server couldn't send jsonData"))
			return
		}
		w.Write([]byte(jsonData))
		log.Printf("Responded to request: order_id %s\n", val)
	}
}
