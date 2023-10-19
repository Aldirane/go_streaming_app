package postgres

import (
	"database/sql"
	"fmt"
	"go_app/pkg/order"
)

func SelectOrder(order_uid string, db *sql.DB) (*order.Order, error) {
	var pass string
	row := db.QueryRow(`select * from orders 
		join delivery on orders.track_number=delivery.track_number 
		join payment on orders.order_uid=payment.transaction where order_uid=$1`, order_uid)
	newOrder := new(order.Order)
	newDelivery := order.Delivery{}
	newPayment := order.Payment{}
	err := row.Scan(
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
		&pass,             // internal database field delivery.track_number
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
		return nil, err
	}
	newOrder.Delivery = newDelivery
	newOrder.Payment = newPayment
	selectItems(newOrder, db)
	return newOrder, nil
}
