package stansub

import (
	"encoding/json"
	"go_app/pkg/order"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
)

var (
	clusterID          string
	clientID           string
	durable            string
	URL                string
	subject            string
	streamCh           = make(chan []byte)
	JsonData           = make(chan *order.Order)
	timeSleepToReceive time.Duration
)

func SubStart() (stan.Conn, stan.Subscription) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	clusterID = os.Getenv("CLUSTER_ID")
	clientID = os.Getenv("CLIENT_SUB_ID")
	durable = os.Getenv("DURABLE")
	URL = os.Getenv("URL")
	subject = os.Getenv("SUBJECT")
	timeSleepToReceive = ParseDuration()
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
	for msg := range streamCh {
		data := new(order.Order)
		err := json.Unmarshal(msg, data)
		if err != nil {
			log.Printf("Wrong data %v\n", err)
		} else {
			JsonData <- data
		}
		time.Sleep(timeSleepToReceive)
	}
}

func ParseDuration() time.Duration {
	envSubSleep := os.Getenv("SUB_SLEEP")
	t, err := strconv.Atoi(envSubSleep)
	if err != nil {
		log.Fatalln("env variable SUB_SLEEP must be integer")
	}
	timeSleep := time.Duration(t * int(time.Second))
	return timeSleep
}
