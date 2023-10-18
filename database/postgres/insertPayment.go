package postgres

import (
	"database/sql"
	"go_app/order"
	"log"
)

func InsertPayment(newOrder *order.Order, db *sql.DB) {
	sqlStrPayment := `INSERT INTO payment(
		transaction, request_id, currency, provider, amount, 
		payment_dt, bank, delivery_cost, goods_total, custom_fee) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT DO NOTHING`
	stmtPayment, err := db.Prepare(sqlStrPayment)
	if err != nil {
		log.Fatalf("Error building statement for payment table %v", err)
	}
	_, err = stmtPayment.Exec(
		newOrder.Payment.Transaction,
		newOrder.Payment.RequestID,
		newOrder.Payment.Currency,
		newOrder.Payment.Provider,
		newOrder.Payment.Amount,
		newOrder.Payment.PaymentDt,
		newOrder.Payment.Bank,
		newOrder.Payment.DeliveryCost,
		newOrder.Payment.GoodsTotal,
		newOrder.Payment.CustomFee,
	)

	if err != nil {
		log.Fatalf("Error cant execute payment statement %v", err)
	}
}
