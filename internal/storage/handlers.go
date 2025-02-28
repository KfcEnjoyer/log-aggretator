package storage

import (
	"auth-service/internal/services"
	"auth-service/internal/user"
	"auth-service/pkg/validator"
	"encoding/json"
	"log"
	"net/http"
)

var v = validator.NewValidator()

func Create(w http.ResponseWriter, r *http.Request, auth *services.AuthService) {
	defer r.Body.Close()

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)

		var u user.RegularUser

		err := decoder.Decode(&u)
		if err != nil {
			panic(err)
		}

		response, err := auth.Register(u)
		if err != nil {
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/jsonF")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func LogIn(w http.ResponseWriter, r *http.Request, auth *services.AuthService) {
	defer r.Body.Close()

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)

		var u user.RegularUser

		err := decoder.Decode(&u)
		if err != nil {
			panic(err)
		}

		response, ok, err := auth.LogIn(u.Username, u.Password.Plain)

		if err != nil || !ok {
			log.Println(err)
			http.Redirect(w, r, "/open", http.StatusFound)
			return
		}

		w.Header().Set("Content-Type", "application/jsonF")
		json.NewEncoder(w).Encode(response)

	}
}
