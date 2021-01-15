package main

import (
	"fmt"
	"log"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

const (
	ServerIp   = "172.16.230.32"
	ServerPort = ":9092"
)

func StartConsumer() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": ServerIp + ServerPort,
		"group.id":          "mygo",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Encounter %s when create kafka consumer", err)
	}

	c.SubscribeTopics([]string{"universe"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
	c.Close()
}
