package services

import (
	"auth-service/internal/database/log"
	"auth-service/internal/database/user_db"
	"auth-service/internal/user"
	"auth-service/pkg/logger"
	"time"
)

type AuthService struct {
	UserRepo *user_db.UserRepoService
	LogRepo  *log.LogRepoService
}

func (a *AuthService) Register(u user.RegularUser) (map[string]string, error) {
	u.Password.Hash()

	if err := a.UserRepo.CreateUser(u); err != nil {
		return nil, err
	}

	var l logger.Log

	response := map[string]string{
		"message":  "User created successfully!",
		"username": u.Username,
		"role":     u.Role,
	}

	l.CreateLog(time.Now(), "INFO", response["message"], u.Username, u.Role)
	if err := a.LogRepo.CreateLog(l); err != nil {
		return nil, err
	}

	return response, nil
}

func (a *AuthService) LogIn(username, password string) (map[string]string, bool, error) {
	var u user.RegularUser
	var l logger.Log

	u, err := a.UserRepo.GetUserByUsername(username)
	if err != nil {
		response := map[string]string{
			"message":  "Error fetching the user",
			"username": u.Username,
			"role":     u.Role,
		}
		return response, false, err
	}

	if !u.Password.Compare(password) {
		response := map[string]string{
			"message":  "Passwords do not match",
			"username": u.Username,
			"role":     u.Role,
		}
		return response, false, err
	}

	response := map[string]string{
		"message":  "User logged in successfully!",
		"username": u.Username,
		"role":     u.Role,
	}

	l.CreateLog(time.Now(), "INFO", response["message"], u.Username, u.Role)
	if err := a.LogRepo.CreateLog(l); err != nil {
		return nil, false, err
	}

	return response, true, nil
}
