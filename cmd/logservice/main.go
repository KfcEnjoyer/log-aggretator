package main

import (
	"auth-service/internal/database"
	log_db "auth-service/internal/database"
	"auth-service/kafka"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs/config.yaml"
	}

	database.CreateConnection(configPath)

	conn := database.GetConnection()
	if conn == nil {
		log.Fatal("Failed to connect to the database")
	}
	defer conn.Close()

	logRepo := log_db.NewLogRepoService()

	if err := logRepo.CreateTable(); err != nil {
		log.Printf("Warning: Failed to create/verify logs table: %v", err)
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	brokers := []string{"localhost:9092"}
	if kafkaBrokers != "" {
		brokers = []string{kafkaBrokers}
	}

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "auth-logs"
	}

	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		groupID = "log-consumer-group"
	}

	consumer := kafka.NewLogConsumer(brokers, topic, groupID, logRepo)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Printf("Starting log consumer service, reading from topic: %s", topic)
	if err := consumer.Start(ctx); err != nil && err != context.Canceled {
		log.Fatalf("Consumer error: %v", err)
	}

	log.Println("Log consumer service stopped")
}
