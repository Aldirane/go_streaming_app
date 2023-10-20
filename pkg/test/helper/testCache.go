package helper

import (
	"fmt"
	"go_app/pkg/cache"
	"go_app/pkg/order"
	"log"

	"github.com/joho/godotenv"
)

// initialize new cache, set new order, get order by id and print it
func CacheTest() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	defaultExpiration, cleanupInterval := cache.ParseCleanupExpiration()
	var orders []*order.Order
	newCache := cache.New(defaultExpiration, cleanupInterval)
	for i := 0; i < 3; i++ {
		orders = append(orders, order.GenerateFakeData())
	}
	testOrder := order.GenerateFakeData()
	newCache.Set(testOrder.OrderID, testOrder, 0)
	if _, ok := newCache.Get(testOrder.OrderID); !ok {
		fmt.Printf("Couldn't get order id %s\n", testOrder.OrderID)
	}
	chachedOrder, _ := newCache.Get(testOrder.OrderID)
	fmt.Printf("%+v\n\n", chachedOrder)
	newCache.SetAllOrders(orders, 0)
	chachedOrders, _ := newCache.GetAllOrders()
	for _, ord := range chachedOrders {
		fmt.Printf("\n%+v\n", ord)
	}
}
