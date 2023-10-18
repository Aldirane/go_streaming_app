package order

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bxcodec/faker/v3"
)

type Order struct {
	OrderID           string   `json:"order_uid" faker:"uuid_digit"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry" faker:"oneof: WBIL"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale" faker:"oneof: ru, en"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id" faker:"username"`
	DeliveryService   string   `json:"delivery_service"`
	ShardKey          string   `json:"shardkey"`
	SmID              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created" faker:"timestamp"`
	OofShard          string   `json:"oof_shard" faker:"oneof: 0, 1"`
}

type Delivery struct {
	Name    string `json:"name" faker:"name"`
	Phone   string `json:"phone" faker:"phone_number"`
	Zip     string `json:"zip" faker:"oneof: 333555, 777653, 555433, 222132, 156258"`
	City    string `json:"city" faker:"oneof: Moscow, St.Peterburg, Omsk, London, Hong-Kong, Berlin"`
	Address string `json:"address" faker:"oneof: Ploshad Mira 15, Lenina 155, Hauptstrasse 21, Schulstrasse 42, Station Road 33, High Street 7"`
	Region  string `json:"region" faker:"oneof: Russia, China, Germany, England, Italy, Spain, Denmark"`
	Email   string `json:"email" faker:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency" faker:"currency"`
	Provider     string `json:"provider" faker:"oneof: wbpay, sberpay, unionpay, tinkoffpay, money order"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt" faker:"unix_time"`
	Bank         string `json:"bank" faker:"oneof: alpha, sber, tinkoff, vtb24"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id" faker:"boundary_start=1, boundary_end=10000000"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name" faker:"username"`
	Sale        int    `json:"sale"`
	Size        string `json:"size" faker:"oneof: 0, XS, XL, L, M, XXL, XML"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand" faker:"oneof: Adidas, Abibas, Nike, Puma, Sumsung, Apple"`
	Status      int    `json:"status"`
}

var (
	minSize = 1
	maxSize = 10
)

func PrintFakeData() {
	order := GenerateFakeData()
	data, err := json.MarshalIndent(order, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", data)

}

func GenerateFakeData() *Order {
	order := new(Order)
	_ = faker.SetRandomMapAndSliceMinSize(minSize)
	_ = faker.SetRandomMapAndSliceMaxSize(maxSize)
	err := faker.FakeData(order)
	if err != nil {
		log.Printf("Cant fake order %v\n", err)
		return nil
	}
	order.Payment.Transaction = order.OrderID
	order_time, err := time.Parse("2006-01-02 15:04:05", order.DateCreated)
	if err != nil {
		log.Printf("Cant parse time %v\n", err)
		return nil
	}
	order.Payment.PaymentDt = order_time.Unix()
	for i := 0; i < len(order.Items); i++ {
		order.Items[i].TrackNumber = order.TrackNumber
	}
	if err != nil {
		fmt.Printf("Cant fake Order %v", err)
	}
	return order
}
