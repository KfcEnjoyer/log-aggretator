package storage

import (
	"auth-service/internal/user"
	"encoding/json"
	"io"
	"net/http"
)

type Storage struct {
	Users map[int]*user.User
}

func (s *Storage) Create(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	if r.Method == "POST" {
		content, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}

		u := new(user.User)

		if err := json.Unmarshal(content, &u); err != nil {
			return err
		}

	}

	return nil
}
