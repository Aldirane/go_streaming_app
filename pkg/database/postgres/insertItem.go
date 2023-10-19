package postgres

import (
	"database/sql"
	"fmt"
	"go_app/pkg/order"
	"log"
	"strings"
)

func InsertItem(newOrder *order.Order, db *sql.DB) {
	sqlStrItem := `INSERT INTO item(
		chrt_id, track_number, price, rid, name, sale, 
		size, total_price, nm_id, brand, status) VALUES`

	rows := []string{}
	itemVals := []interface{}{}

	for idx, item := range newOrder.Items {
		fmtRow := []string{}
		for i := 1; i <= 11; i++ {
			fmtRow = append(fmtRow, fmt.Sprintf("$%d", idx*11+i))
		}
		finalRow := "(" + strings.Join(fmtRow, ", ") + ")"
		rows = append(rows, finalRow)
		itemVals = append(
			itemVals,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
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
