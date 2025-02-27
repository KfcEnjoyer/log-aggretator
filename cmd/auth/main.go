package main

import (
	"auth-service/internal/database"
	database2 "auth-service/internal/database/log"
	database1 "auth-service/internal/database/user_db"
	"auth-service/internal/storage"
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	conn := database.GetConnection()
	if conn == nil {
		log.Fatal("Failed to connect to the database")
	}
	defer conn.Close(context.Background())

	database1.CreateTable()
	database2.CreateTable()

	port := "localhost:8080"
	http.HandleFunc("/create", storage.Create)
	http.HandleFunc("/login", storage.LogIn)
	http.HandleFunc("/home", Home)

	http.ListenAndServe(port, nil)

}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello!")
}
