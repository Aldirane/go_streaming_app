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

func DbConnect(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Cant open connection: %s\n%v\n", connStr, err)
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
