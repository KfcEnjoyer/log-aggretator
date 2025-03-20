package main

import (
	"auth-service/internal/database"
	"auth-service/internal/services"
	"auth-service/pkg/logger"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	VERSION     = "1.0.0"
	ENVIRONMENT = "development"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "../configs/config.yaml"
	}

	database.CreateConnection(configPath)

	conn := database.GetConnection()
	if conn == nil {
		log.Fatal("Failed to connect to the database")
	}
	defer conn.Close()

	userRepo := database.NewUserRepoService()
	logRepo := database.NewLogRepoService()

	if err := logRepo.CreateTable(); err != nil {
		log.Printf("Warning: Failed to create/verify logs table: %v", err)
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	brokers := []string{"localhost:9092"}
	if kafkaBrokers != "" {
		brokers = []string{kafkaBrokers}
	}

	auth := services.NewAuthService(userRepo, logRepo, brokers)
	defer auth.Close()

	var l logger.Log
	l.SystemStartup(time.Now(), VERSION, ENVIRONMENT)
	if auth.KafkaProducer != nil {
		auth.KafkaProducer.SendLog(l)
	} else {
		auth.LogRepo.CreateLog(l)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	serverAddr := fmt.Sprintf("0.0.0.0:%s", port)

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		services.Create(w, r, auth)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		services.LogIn(w, r, auth)
	})
	http.HandleFunc("/home", Home)
	http.HandleFunc("/open", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Open endpoint accessed")
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Auth service is running. Version: %s", VERSION)
	})

	server := &http.Server{
		Addr: serverAddr,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		server.Close()
	}()

	log.Printf("Starting auth service on %s...", serverAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the authentication service!")
}
