package helper

import (
	"fmt"
	"go_app/pkg/database"
	"go_app/pkg/database/postgres"
	"go_app/pkg/order"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Connect to database
func DatabaseTest() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var (
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
		sslmode  = os.Getenv("SSLMODE")
		ConnStr  = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, dbname, sslmode)
	)
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
