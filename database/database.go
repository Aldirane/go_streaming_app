package database

import (
	"database/sql"
	"fmt"
	"go_app/order"
	"log"
	"strings"

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
	ord, err := SelectOrder(newOrder.OrderID, db)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", ord)
	newOrders, err := SelectOrders(db)
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
	InsertOrder(newOrder, db)
	InsertDelivery(newOrder, db)
	InsertPayment(newOrder, db)
	InsertItem(newOrder, db)
}

func selectItems(newOrder *order.Order, db *sql.DB) {
	var items []order.Item
	rows, err := db.Query("select * from item where item.track_number=$1", newOrder.TrackNumber)
	if err != nil {
		log.Printf("Wrong query statement for item table %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		item := order.Item{}
		err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price,
			&item.RID, &item.Name, &item.Sale, &item.Size, &item.TotalPrice,
			&item.NmID, &item.Brand, &item.Status)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, item)
	}
	newOrder.Items = items
}

func SelectOrders(db *sql.DB) ([]*order.Order, error) {
	var pass string
	var newOrders []*order.Order
	queryStr := "select * from orders " +
		"join delivery on orders.track_number=delivery.track_number " +
		"join payment on orders.order_uid=payment.transaction"
	queryStmt, err := db.Prepare(queryStr)
	if err != nil {
		log.Printf("Wrong query string %v", err)
		return nil, err
	}
	rows, err := queryStmt.Query()
	if err != nil {
		log.Printf("Wrong query statement %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		newOrder := new(order.Order)
		newDelivery := order.Delivery{}
		newPayment := order.Payment{}
		err := rows.Scan(
			&newOrder.OrderID, &newOrder.TrackNumber, &newOrder.Entry, &newOrder.Locale,
			&newOrder.InternalSignature, &newOrder.CustomerID, &newOrder.DeliveryService,
			&newOrder.ShardKey, &newOrder.SmID, &newOrder.DateCreated, &newOrder.OofShard,
			&pass, &newDelivery.Name, &newDelivery.Phone, &newDelivery.Zip, &newDelivery.City,
			&newDelivery.Address, &newDelivery.Region, &newDelivery.Email,
			&newPayment.Transaction, &newPayment.RequestID, &newPayment.Currency, &newPayment.Provider,
			&newPayment.Amount, &newPayment.PaymentDt, &newPayment.Bank,
			&newPayment.DeliveryCost, &newPayment.GoodsTotal, &newPayment.CustomFee,
		)
		if err != nil {
			fmt.Println(err)
			continue
		}
		newOrder.Delivery = newDelivery
		newOrder.Payment = newPayment
		selectItems(newOrder, db)
		newOrders = append(newOrders, newOrder)
	}
	return newOrders, nil
}

func SelectOrder(order_uid string, db *sql.DB) (*order.Order, error) {
	var pass string
	row := db.QueryRow("select * from orders "+
		"join delivery on orders.track_number=delivery.track_number "+
		"join payment on orders.order_uid=payment.transaction where order_uid=$1", order_uid)
	newOrder := new(order.Order)
	newDelivery := order.Delivery{}
	newPayment := order.Payment{}
	err := row.Scan(
		&newOrder.OrderID, &newOrder.TrackNumber, &newOrder.Entry, &newOrder.Locale,
		&newOrder.InternalSignature, &newOrder.CustomerID, &newOrder.DeliveryService,
		&newOrder.ShardKey, &newOrder.SmID, &newOrder.DateCreated, &newOrder.OofShard,
		&pass, &newDelivery.Name, &newDelivery.Phone, &newDelivery.Zip, &newDelivery.City,
		&newDelivery.Address, &newDelivery.Region, &newDelivery.Email,
		&newPayment.Transaction, &newPayment.RequestID, &newPayment.Currency, &newPayment.Provider,
		&newPayment.Amount, &newPayment.PaymentDt, &newPayment.Bank,
		&newPayment.DeliveryCost, &newPayment.GoodsTotal, &newPayment.CustomFee,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	newOrder.Delivery = newDelivery
	newOrder.Payment = newPayment
	selectItems(newOrder, db)
	return newOrder, nil
}

func InsertOrder(newOrder *order.Order, db *sql.DB) {
	sqlStrOrder := "INSERT INTO orders(" +
		"order_uid, track_number, entry, " +
		"locale, internal_signature, customer_id, " +
		"delivery_service, shardkey, sm_id, date_created, oof_shard) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT DO NOTHING"
	stmtOrder, err := db.Prepare(sqlStrOrder)
	if err != nil {
		log.Fatalf("Error building statement for orders table %v", err)
	}
	_, err = stmtOrder.Exec(
		newOrder.OrderID, newOrder.TrackNumber, newOrder.Entry, newOrder.Locale,
		newOrder.InternalSignature, newOrder.CustomerID, newOrder.DeliveryService, newOrder.ShardKey,
		newOrder.SmID, newOrder.DateCreated, newOrder.OofShard)

	if err != nil {
		log.Fatalf("Error cant execute order statement %v", err)
	}
}

func InsertDelivery(newOrder *order.Order, db *sql.DB) {
	sqlStrDelivery := "INSERT INTO delivery(" +
		"track_number, name, phone, zip, city, address, region, email) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT DO NOTHING"
	stmtDelivery, err := db.Prepare(sqlStrDelivery)
	if err != nil {
		log.Fatalf("Error building statement for delivery table %v", err)
	}
	_, err = stmtDelivery.Exec(
		newOrder.TrackNumber, newOrder.Delivery.Name, newOrder.Delivery.Phone, newOrder.Delivery.Zip,
		newOrder.Delivery.City, newOrder.Delivery.Address, newOrder.Delivery.Region, newOrder.Delivery.Email)

	if err != nil {
		log.Fatalf("Error cant execute delivery statement %v", err)
	}
}

func InsertPayment(newOrder *order.Order, db *sql.DB) {
	sqlStrPayment := "INSERT INTO payment(" +
		"transaction, request_id, currency, provider, amount, " +
		"payment_dt, bank, delivery_cost, goods_total, custom_fee) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT DO NOTHING"
	stmtPayment, err := db.Prepare(sqlStrPayment)
	if err != nil {
		log.Fatalf("Error building statement for payment table %v", err)
	}
	_, err = stmtPayment.Exec(
		newOrder.Payment.Transaction, newOrder.Payment.RequestID, newOrder.Payment.Currency,
		newOrder.Payment.Provider, newOrder.Payment.Amount, newOrder.Payment.PaymentDt, newOrder.Payment.Bank,
		newOrder.Payment.DeliveryCost, newOrder.Payment.GoodsTotal, newOrder.Payment.CustomFee)

	if err != nil {
		log.Fatalf("Error cant execute payment statement %v", err)
	}
}

func InsertItem(newOrder *order.Order, db *sql.DB) {
	sqlStrItem := "INSERT INTO item(" +
		"chrt_id, track_number, price, rid, name, sale, " +
		"size, total_price, nm_id, brand, status) VALUES"

	rows := []string{}
	itemVals := []interface{}{}

	for idx, item := range newOrder.Items {
		fmtRow := []string{}
		for i := 1; i <= 11; i++ {
			fmtRow = append(fmtRow, fmt.Sprintf("$%d", idx*11+i))
		}
		finalRow := "(" + strings.Join(fmtRow, ", ") + ")"
		rows = append(rows, finalRow)
		itemVals = append(itemVals, item.ChrtID, item.TrackNumber, item.Price, item.RID,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
	}
	sqlStrItem += strings.Join(rows, " , ") + " ON CONFLICT DO NOTHING"
	stmtItem, err := db.Prepare(sqlStrItem)
	if err != nil {
		log.Fatalf("Error building statement for item table %v", err)
	}
	_, err = stmtItem.Exec(itemVals...)

	if err != nil {
		log.Fatalf("Error cant execute item statement %v", err)
	}
}
