package postgres

import (
	"database/sql"
	"fmt"
	"go_app/pkg/order"
	"log"
)

func SelectOrders(db *sql.DB, orderBy string, ascDesc string, limit int, offset int) ([]*order.Order, error) {
	var (
		limitRows                interface{}
		RowsOrderBy, RowsAscDesc string
		pass                     string
		newOrders                []*order.Order
	)
	limitRows = limit
	if limit == 0 {
		limitRows = "all"
	}
	if orderBy == "" {
		RowsOrderBy = "date_created"
	}
	if ascDesc == "" {
		RowsAscDesc = "desc"
	}

	queryStr := fmt.Sprintf(`select * from orders 
		join delivery on orders.track_number=delivery.track_number 
		join payment on orders.order_uid=payment.transaction 
		order by %s %s limit %v offset %v`, RowsOrderBy, RowsAscDesc, limitRows, offset)
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
