package postgres

import (
	"database/sql"
	"fmt"
	"go_app/order"
	"log"
)

func SelectOrders(db *sql.DB) ([]*order.Order, error) {
	var pass string
	var newOrders []*order.Order
	queryStr := `select * from orders 
		join delivery on orders.track_number=delivery.track_number 
		join payment on orders.order_uid=payment.transaction`
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
			&newOrder.OrderID, // Order data
			&newOrder.TrackNumber,
			&newOrder.Entry,
			&newOrder.Locale,
			&newOrder.InternalSignature,
			&newOrder.CustomerID,
			&newOrder.DeliveryService,
			&newOrder.ShardKey,
			&newOrder.SmID,
			&newOrder.DateCreated,
			&newOrder.OofShard,
			&pass,             // Internal database field delivery.track_number
			&newDelivery.Name, // Delivery data
			&newDelivery.Phone,
			&newDelivery.Zip,
			&newDelivery.City,
			&newDelivery.Address,
			&newDelivery.Region,
			&newDelivery.Email,
			&newPayment.Transaction, // Payment data
			&newPayment.RequestID,
			&newPayment.Currency,
			&newPayment.Provider,
			&newPayment.Amount,
			&newPayment.PaymentDt,
			&newPayment.Bank,
			&newPayment.DeliveryCost,
			&newPayment.GoodsTotal,
			&newPayment.CustomFee,
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
