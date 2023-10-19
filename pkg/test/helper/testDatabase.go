package helper

import (
	"fmt"
	"go_app/pkg/database"
	"go_app/pkg/database/postgres"
	"go_app/pkg/order"
	"log"

	_ "github.com/lib/pq"
)

var (
	user     = "aldar"
	password = "password"
	dbname   = "aldar"
	sslmode  = "disable"
	ConnStr  = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, dbname, sslmode)
)

// Connect to database
func DatabaseTest() {
	newOrder := order.GenerateFakeData()
	db := database.DbConnect(ConnStr)
	defer db.Close()
	database.InsertOrderAndAll(newOrder, db)
	ord, err := postgres.SelectOrder(newOrder.OrderID, db)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", ord)
	newOrders, err := postgres.SelectOrders(db)
	if err != nil {
		log.Fatalln(err)
	}
	for _, ord := range newOrders {
		fmt.Printf("%+v\n", ord)
	}
}
