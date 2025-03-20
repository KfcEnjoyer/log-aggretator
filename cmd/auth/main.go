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
	database.CreateConnection("C:\\Users\\user\\Desktop\\coding and stuff\\study or portfolio projects\\auth-service\\configs\\config.yaml")

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

	brokers := []string{"localhost:9092"}

	auth := services.NewAuthService(userRepo, logRepo, brokers)
	defer auth.Close()

	var l logger.Log
	l.SystemStartup(time.Now(), VERSION, ENVIRONMENT)
	if auth.KafkaProducer != nil {
		auth.KafkaProducer.SendLog(l)
	} else {
		auth.LogRepo.CreateLog(l)
	}

	port := "localhost:8080"
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		services.Create(w, r, auth)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		services.LogIn(w, r, auth)
	})
	//http.HandleFunc("/update-role", func(w http.ResponseWriter, r *http.Request) {
	//	services.UpdateRole(w, r, auth)
	//})
	http.HandleFunc("/home", Home)
	http.HandleFunc("/open", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Open endpoint accessed")
	})

	server := &http.Server{
		Addr: port,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		server.Close()
	}()

	log.Printf("Starting auth service on %s...", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the authentication service!")
}
