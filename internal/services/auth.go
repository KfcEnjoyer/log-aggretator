package services

import (
	"auth-service/internal/database"
	"auth-service/internal/user"
	"auth-service/kafka"
	"auth-service/pkg/logger"
	stdlog "log"
	"time"
)

type AuthService struct {
	UserRepo      *database.UserRepoService
	LogRepo       *database.LogRepoService
	KafkaProducer *kafka.LogProducer
}

func NewAuthService(userRepo *database.UserRepoService, logRepo *database.LogRepoService, brokers []string) *AuthService {
	producer := kafka.NewLogProducer(brokers, "auth-logs")

	return &AuthService{
		UserRepo:      userRepo,
		LogRepo:       logRepo,
		KafkaProducer: producer,
	}
}

func (a *AuthService) Close() {
	if a.KafkaProducer != nil {
		a.KafkaProducer.Close()
	}
}

func (a *AuthService) Register(u user.RegularUser, ipAddress string) (map[string]string, error) {
	u.Password.Hash()

	if err := a.UserRepo.CreateUser(u); err != nil {
		var l logger.Log
		l.DetailedLog(
			time.Now(),
			logger.ERROR,
			"Failed to create user account",
			logger.ACCOUNT,
			u.Username,
			u.Role,
			map[string]interface{}{
				"error":      err.Error(),
				"ip_address": ipAddress,
			},
		)

		if a.KafkaProducer != nil {
			if err := a.KafkaProducer.SendLog(l); err != nil {
				stdlog.Printf("Failed to send log to Kafka: %v", err)
			}
		} else {
			if err := a.LogRepo.CreateLog(l); err != nil {
				stdlog.Printf("Failed to create log: %v", err)
			}
		}

		return nil, err
	}

	var l logger.Log
	l.UserCreated(time.Now(), u.Username, u.Role, ipAddress)

	if a.KafkaProducer != nil {
		if err := a.KafkaProducer.SendLog(l); err != nil {
			stdlog.Printf("Failed to send log to Kafka: %v", err)
		}
	} else {
		if err := a.LogRepo.CreateLog(l); err != nil {
			stdlog.Printf("Failed to create log: %v", err)
		}
	}

	response := map[string]string{
		"message":  "User created successfully!",
		"username": u.Username,
		"role":     u.Role,
	}

	return response, nil
}

func (a *AuthService) LogIn(username, password, ipAddress, userAgent string) (map[string]string, bool, error) {
	var l logger.Log

	u, err := a.UserRepo.GetUserByUsername(username)
	if err != nil {
		l.UserNotFound(time.Now(), username, ipAddress)

		if a.KafkaProducer != nil {
			a.KafkaProducer.SendLog(l)
		} else {
			a.LogRepo.CreateLog(l)
		}

		response := map[string]string{
			"message":  "User not found",
			"username": username,
		}
		return response, false, err
	}

	if !u.Password.Compare(password) {
		l.LoginFailure(time.Now(), u.Username, u.Role, "incorrect_password", ipAddress)

		if a.KafkaProducer != nil {
			a.KafkaProducer.SendLog(l)
		} else {
			a.LogRepo.CreateLog(l)
		}

		response := map[string]string{
			"message":  "Incorrect password",
			"username": u.Username,
			"role":     u.Role,
		}
		return response, false, nil
	}

	l.LoginSuccess(time.Now(), u.Username, u.Role, ipAddress)

	if a.KafkaProducer != nil {
		a.KafkaProducer.SendLog(l)
	} else {
		a.LogRepo.CreateLog(l)
	}

	response := map[string]string{
		"message":  "User logged in successfully!",
		"username": u.Username,
		"role":     u.Role,
	}

	return response, true, nil
}

func (a *AuthService) LogActivity(username, role, activity, ipAddress string) error {
	var l logger.Log
	l.DetailedLog(
		time.Now(),
		logger.INFO,
		activity,
		"activity",
		username,
		role,
		map[string]interface{}{
			"ip_address": ipAddress,
			"timestamp":  time.Now().Format(time.RFC3339),
		},
	)

	if a.KafkaProducer != nil {
		return a.KafkaProducer.SendLog(l)
	}

	return a.LogRepo.CreateLog(l)
}
