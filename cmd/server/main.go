package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"go_app/cache"
	"go_app/database"
	"go_app/stansub"
)

var (
	db                = database.DbConnect()
	defaultExpiration = 30 * time.Minute
	cleanupInterval   = 40 * time.Minute
	newCache          = cache.New(defaultExpiration, cleanupInterval)
	signalChan        = make(chan os.Signal, 1)
)

func listenSub(wg *sync.WaitGroup) {
	defer wg.Done()
	for order := range stansub.JsonData {
		log.Printf("Received new order: ID = %s\n", order.OrderID)
		newCache.Set(order.OrderID, order, 0)
		database.InsertOrderAndAll(&order, db)
	}
}

// First ctrl+c stopping stan subscription, second command ctrl+c stop server

func main() {
	wg := new(sync.WaitGroup)
	sc, sub := stansub.SubStart()
	wg.Add(1)
	go listenSub(wg)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for range signalChan {
			close(signalChan)
			signal.Stop(signalChan)
			stansub.SubClose(sc, sub)
		}
	}(wg)
	signal.Notify(signalChan, os.Interrupt)
	orders, err := database.SelectOrders(db)
	if err != nil {
		log.Fatal(err)
	}
	newCache.SetAllOrders(orders, 0)
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
