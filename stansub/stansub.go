package stansub

import (
	"encoding/json"
	"go_app/order"
	"log"
	"time"

	"github.com/nats-io/stan.go"
)

var (
	clusterID          = "test-cluster"
	clientID           = "client-sub-order"
	durable            = "client-sub-order"
	URL                = "nats://0.0.0.0:4222"
	subject            = "orders"
	streamCh           = make(chan []byte)
	JsonData           = make(chan order.Order)
	timeSleepToReceive = time.Second * 2
)

func SubStart() (stan.Conn, stan.Subscription) {
	sc, sub := SubscribtionConnect()
	go unmarshalStreamData()
	return sc, sub
}

func SubscribtionConnect() (stan.Conn, stan.Subscription) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(URL),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Printf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Printf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s\n", err, URL)
		return nil, nil
	}
	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", URL, clusterID, clientID)
	sub, err := sc.Subscribe(subject, func(m *stan.Msg) {
		streamCh <- m.Data
	}, stan.DurableName("client-sub-durable"))
	if err != nil {
		sc.Close()
		log.Println(err)
	}
	log.Printf("Listening on [%s], clientID=[%s], durable=[%s]\n", subject, clientID, durable)
	return sc, sub
}

func SubClose(sc stan.Conn, sub stan.Subscription) {
	log.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
	if durable == "" {
		sub.Unsubscribe()
	}
	sc.Close()
}

func unmarshalStreamData() {
	var data order.Order
	for msg := range streamCh {
		err := json.Unmarshal(msg, &data)
		if err != nil {
			log.Printf("Wrong data %v\n", err)
		} else {
			JsonData <- data
		}
		time.Sleep(timeSleepToReceive)
	}
}
