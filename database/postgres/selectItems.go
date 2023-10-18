package postgres

import (
	"database/sql"
	"fmt"
	"go_app/order"
	"log"
)

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
