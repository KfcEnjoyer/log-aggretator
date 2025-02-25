package main

import (
	"auth-service/internal/database"
	"auth-service/internal/storage"
	"context"
	"log"
	"net/http"
)

func main() {
	conn := database.CreateConnection()
	if conn == nil {
		log.Fatal("Failed to connect to the database")
	}
	defer conn.Close(context.Background())

	database.CreateTable(conn)

	port := "localhost:8080"
	http.HandleFunc("/create", storage.Create)

	http.ListenAndServe(port, nil)

}
