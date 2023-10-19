package main

import (
	"encoding/json"
	"log"
	"time"

	"go_app/pkg/order"

	"github.com/nats-io/stan.go"
)

func main() {
	var (
		clusterID       = "test-cluster"
		clientID        = "client-pub"
		URL             = "nats://0.0.0.0:4222"
		subject         = "orders"
		timeSleepToSend = 10 * time.Second
	)
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(URL))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	defer sc.Close()

	for i := 0; ; i++ {
		var fake_order interface{}
		if i%5 == 0 || i%3 == 0 {
			fake_order = "Wrong data"
		} else {
			fake_order = order.GenerateFakeData()
		}
		msg, err := json.Marshal(fake_order)
		if err != nil {
			log.Printf("Error during marshaling %v\n", err)
		}
		err = sc.Publish(subject, msg)
		if err != nil {
			log.Fatalf("Error during publish: %v\n", err)
		}
		log.Printf("Published [%s] : '%d'\n", subject, i)
		time.Sleep(timeSleepToSend)
	}
}
