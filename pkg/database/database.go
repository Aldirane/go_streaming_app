package database

import (
	"database/sql"
	"fmt"
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
)

var ConnStr = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, dbname, sslmode)

func DatabaseTest() {
	newOrder := order.GenerateFakeData()
	db := DbConnect()
	defer db.Close()
	InsertOrderAndAll(newOrder, db)
	ord, err := postgres.SelectOrder(newOrder.OrderID, db)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", ord)
	newOrders, err := postgres.SelectOrders(db)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", newOrders)
	for _, order := range newOrders {
		fmt.Printf("%+v\n", order)
	}
}

func DbConnect() *sql.DB {
	db, err := sql.Open("postgres", ConnStr)
	if err != nil {
		log.Fatalf("Cant open connection: %s\n%v\n", ConnStr, err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("can't ping database %v\n", err)
	} else {
		log.Println("Database connected successfuly")
	}
	return db
}

func InsertOrderAndAll(newOrder *order.Order, db *sql.DB) {
	postgres.InsertOrder(newOrder, db)
	postgres.InsertDelivery(newOrder, db)
	postgres.InsertPayment(newOrder, db)
	postgres.InsertItem(newOrder, db)
}
