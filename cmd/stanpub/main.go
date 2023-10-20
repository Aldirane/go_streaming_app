package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"go_app/pkg/order"

	"github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	clusterID := os.Getenv("CLUSTER_ID")
	clientID := os.Getenv("CLIENT_PUB_ID")
	URL := os.Getenv("URL")
	subject := os.Getenv("SUBJECT")
	timeSleepToSend := ParseDuration()
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

func ParseDuration() time.Duration {
	envSubSleep := os.Getenv("PUB_SLEEP")
	t, err := strconv.Atoi(envSubSleep)
	if err != nil {
		log.Fatalln("env variable PUB_SLEEP must be integer")
	}
	timeSleep := time.Duration(t * int(time.Second))
	return timeSleep
}
