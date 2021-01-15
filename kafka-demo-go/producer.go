package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func readFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	tokens := make([]string, 100)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		tokens = append(tokens, strings.Fields(line)...)
	}
	return tokens, nil
}

func parseArgs() (*string, *string) {
	novelPath := flag.String("p", "", "Path to novel")
	topic := flag.String("t", "", "Topic name")

	flag.Parse()

	if *novelPath == "" || *topic == "" {
		log.Fatal("Unable to parse filepath and topic from cmdline")
	}
	return novelPath, topic
}

func CreateTopic(p *kafka.Producer, topic *string) error {
	admin, err := kafka.NewAdminClientFromProducer(p)
	if err != nil {
		log.Printf("Failed to create admin client from producer: %s", err)
		return err
	}
	defer admin.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maxDuration, err := time.ParseDuration("5s")
	if err != nil {
		return err
	}

	results, err := admin.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{{
			Topic:             *topic,
			NumPartitions:     1,
			ReplicationFactor: 3,
		}},
		kafka.SetAdminOperationTimeout(maxDuration),
	)
	if err != nil {
		return err
	}

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			log.Printf("Failed to create topic: %v", result.Error)
			return result.Error
		}
	}
	return nil
}

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "172.16.230.32:9092"})
	// p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "test-group-1-vm1:9092"})
	if err != nil {
		log.Fatalf("Encountering: %s", err)
	}
	defer p.Close()

	path, topic := parseArgs()
	err = CreateTopic(p, topic)
	if err != nil {
		log.Fatalf("Unable to create topic: %s", err)
	}

	tokens, err := readFromFile(*path)
	if err != nil {
		log.Fatalf("Unable to open file: %s", *path)
	}

	for i, word := range tokens {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
			Key:            []byte(strconv.Itoa(i % 5)),
			Value:          []byte(word),
		}, nil)
	}
	p.Flush(10000)
}
