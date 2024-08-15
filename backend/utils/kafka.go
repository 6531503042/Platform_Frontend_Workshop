package utils

import (
	"backend/config"
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

var KafkaWriter *kafka.Writer
var KafkaReader *kafka.Reader

// InitKafka initializes the Kafka writer and reader
func InitKafka() {
	// Initialize Kafka Writer
	KafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(config.KafkaBroker),
		Topic:    config.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}

	// Initialize Kafka Reader
	KafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.KafkaBroker},
		Topic:   config.KafkaTopic,
		GroupID: "test-group",
	})

	// Check Kafka Writer and Reader connection (Optional)
	if err := checkKafkaConnection(); err != nil {
		log.Fatalf("Error initializing Kafka: %v", err)
	}
}

// checkKafkaConnection attempts to write and read a test message to/from Kafka
func checkKafkaConnection() error {
	// Write a test message
	err := KafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("test-key"),
			Value: []byte("test-message"),
		},
	)
	if err != nil {
		return err
	}

	// Read a test message
	_, err = KafkaReader.ReadMessage(context.Background())
	if err != nil {
		return err
	}

	log.Println("Kafka connection successful!")
	return nil
}

// CloseKafka closes Kafka writer and reader resources
func CloseKafka() {
	if err := KafkaWriter.Close(); err != nil {
		log.Printf("Error closing KafkaWriter: %v", err)
	}

	if err := KafkaReader.Close(); err != nil {
		log.Printf("Error closing KafkaReader: %v", err)
	}
}