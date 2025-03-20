package main

import (
	"auth-service/internal/database"
	log_db "auth-service/internal/database"
	"auth-service/kafka"
	"context"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	database.CreateConnection("C:\\Users\\user\\Desktop\\coding and stuff\\study or portfolio projects\\auth-service\\configs\\config.yaml")

	conn := database.GetConnection()
	if conn == nil {
		log.Fatal("Failed to connect to the database")
	}
	defer conn.Close()

	logRepo := log_db.NewLogRepoService()

	if err := logRepo.CreateTable(); err != nil {
		log.Printf("Warning: Failed to create/verify logs table: %v", err)
	}

	brokers := []string{"localhost:9092"}
	topic := "auth-logs"
	groupID := "log-consumer-group"

	consumer := kafka.NewLogConsumer(brokers, topic, groupID, logRepo)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Printf("Starting log consumer service, reading from topic: %s", topic)
	if err := consumer.Start(ctx); err != nil && err != context.Canceled {
		log.Fatalf("Consumer error: %v", err)
	}

	log.Println("Log consumer service stopped")
}
