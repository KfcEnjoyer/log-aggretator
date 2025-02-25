package storage

import (
	"auth-service/internal/database"
	"auth-service/internal/user"
	"encoding/json"
	"net/http"
)

var conn = database.CreateConnection()

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
			"user":    u.Username,
		}

		database.CreateUser(conn, u)

		w.Header().Set("Content-Type", "application/jsonF")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
