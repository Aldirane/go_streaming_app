package postgres

import (
	"database/sql"
	"go_app/pkg/order"
	"log"
)

func InsertOrder(newOrder *order.Order, db *sql.DB) {
	sqlStrOrder := `INSERT INTO orders(
		order_uid, track_number, entry, 
		locale, internal_signature, customer_id, 
		delivery_service, shardkey, sm_id, date_created, oof_shard) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT DO NOTHING`
	stmtOrder, err := db.Prepare(sqlStrOrder)
	if err != nil {
		log.Fatalf("Error building statement for orders table %v", err)
	}
	_, err = stmtOrder.Exec(
		newOrder.OrderID,
		newOrder.TrackNumber,
		newOrder.Entry,
		newOrder.Locale,
		newOrder.InternalSignature,
		newOrder.CustomerID,
		newOrder.DeliveryService,
		newOrder.ShardKey,
		newOrder.SmID,
		newOrder.DateCreated,
		newOrder.OofShard,
	)

	if err != nil {
		log.Fatalf("Error cant execute order statement %v", err)
	}
}
