package postgres

import (
	"database/sql"
	"go_app/order"
	"log"
)

func InsertDelivery(newOrder *order.Order, db *sql.DB) {
	sqlStrDelivery := `INSERT INTO delivery(
		track_number, name, phone, zip, city, address, region, email) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT DO NOTHING`
	stmtDelivery, err := db.Prepare(sqlStrDelivery)
	if err != nil {
		log.Fatalf("Error building statement for delivery table %v", err)
	}
	_, err = stmtDelivery.Exec(
		newOrder.TrackNumber,
		newOrder.Delivery.Name,
		newOrder.Delivery.Phone,
		newOrder.Delivery.Zip,
		newOrder.Delivery.City,
		newOrder.Delivery.Address,
		newOrder.Delivery.Region,
		newOrder.Delivery.Email,
	)

	if err != nil {
		log.Fatalf("Error cant execute delivery statement %v", err)
	}
}
