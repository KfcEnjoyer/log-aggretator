package main

import (
	"auth-service/internal/database"
	"context"
	"log"
)

func main() {
	conn := database.CreateConnection()
	if conn == nil {
		log.Fatal("Failed to connect to the database")
	}
	defer conn.Close(context.Background())

	database.CreateTable(conn)

	database.CreateUser(conn, "Nsdad", "dasdasd")
}
