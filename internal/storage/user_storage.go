package storage

import (
	database2 "auth-service/internal/database/log"
	database1 "auth-service/internal/database/user_db"
	"auth-service/internal/user"
	"auth-service/pkg/logger"
	"auth-service/pkg/validator"
	"encoding/json"
	"net/http"
	"time"
)

var v = validator.NewValidator()

func Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)

		var u user.RegularUser

		err := decoder.Decode(&u)
		if err != nil {
			panic(err)
		}

		u.Password.Hash()
		response := map[string]string{
			"message": "Created successfully!",
			"user_db": u.Username,
		}

		now := time.Now().UTC()

		var l logger.Log
		l.CreateLog(now, "INFO", response["message"], u.Username, u.Role)

		database2.CreateLog(l)
		database1.CreateUser(u)

		w.Header().Set("Content-Type", "application/jsonF")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)

		var u user.RegularUser

		err := decoder.Decode(&u)
		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		if !database1.GetUser(u.Username) {
			http.Error(w, "username not found", http.StatusNotFound)
			return
		}

		p := database1.GetUserPassword(u.Username)

		if !u.Password.Compare(p) {
			http.Error(w, "Passwords do not match", http.StatusUnauthorized)
			return
		}

		response := map[string]string{
			"message": "Logged in successfully!",
			"user_db": u.Username,
		}

		w.Header().Set("Content-Type", "application/json")
		http.Redirect(w, r, "/home", http.StatusFound)
		json.NewEncoder(w).Encode(response)
	}
}
