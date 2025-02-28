package main

import (
	"auth-service/internal/database"
	log2 "auth-service/internal/database/log"
	"auth-service/internal/database/user_db"
	"auth-service/internal/services"
	"auth-service/internal/storage"
	"fmt"
	"log"
	"net/http"
)

func main() {
	database.CreateConnection("C:\\Users\\user\\Desktop\\coding and stuff\\study or portfolio projects\\auth-service\\configs\\config.yaml")

	conn := database.GetConnection()

	if conn == nil {
		log.Fatal("Failed to connect to the database")
	}
	defer conn.Close()

	userRepo := user_db.NewUserRepoService()
	logRepo := log2.NewLogRepoService() // Assuming this function exists
	auth := services.AuthService{
		UserRepo: userRepo,
		LogRepo:  logRepo,
	}

	port := "localhost:8080"
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		storage.Create(w, r, &auth)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		storage.LogIn(w, r, &auth)
	})
	http.HandleFunc("/home", Home)
	http.HandleFunc("/open", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("NFFFF")
	})

	http.ListenAndServe(port, nil)

}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello!")
}
