package logger

import (
	"encoding/json"
	"time"
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

const (
	AUTH_SUCCESS = "auth_success"
	AUTH_FAILURE = "auth_failure"
	ACCOUNT      = "account"
	SECURITY     = "security"
	SYSTEM       = "system"
	PERFORMANCE  = "performance"
)

type Log struct {
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Category  string                 `json:"category"`
	TimeStamp time.Time              `json:"timestamp"`
	Username  string                 `json:"username"`
	Role      string                 `json:"role"`
	Metadata  map[string]interface{} `json:"metadata"`
}

func (l *Log) CreateLog(timestamp time.Time, level, message, username, role string) {
	l.TimeStamp = timestamp
	l.Level = level
	l.Message = message
	l.Username = username
	l.Role = role
	l.Category = "auth"
	l.Metadata = make(map[string]interface{})
}

func (l *Log) DetailedLog(timestamp time.Time, level, message, category, username, role string, metadata map[string]interface{}) {
	l.TimeStamp = timestamp
	l.Level = level
	l.Message = message
	l.Category = category
	l.Username = username
	l.Role = role
	l.Metadata = metadata
}

func (l *Log) LoginSuccess(timestamp time.Time, username, role, ipAddress string) {
	l.DetailedLog(
		timestamp,
		INFO,
		"User logged in successfully",
		AUTH_SUCCESS,
		username,
		role,
		map[string]interface{}{
			"ip_address": ipAddress,
			"login_time": timestamp.Format(time.RFC3339),
		},
	)
}

func (l *Log) LoginFailure(timestamp time.Time, username, role, reason, ipAddress string) {
	l.DetailedLog(
		timestamp,
		WARN,
		"Login attempt failed: "+reason,
		AUTH_FAILURE,
		username,
		role,
		map[string]interface{}{
			"failure_reason": reason,
			"ip_address":     ipAddress,
			"timestamp":      timestamp.Format(time.RFC3339),
		},
	)
}

func (l *Log) UserCreated(timestamp time.Time, username, role, ipAddress string) {
	l.DetailedLog(
		timestamp,
		INFO,
		"New user account created",
		ACCOUNT,
		username,
		role,
		map[string]interface{}{
			"ip_address":        ipAddress,
			"registration_time": timestamp.Format(time.RFC3339),
		},
	)
}

func (l *Log) UserNotFound(timestamp time.Time, username, ipAddress string) {
	l.DetailedLog(
		timestamp,
		WARN,
		"Login attempt for non-existent user",
		SECURITY,
		username,
		"unknown",
		map[string]interface{}{
			"ip_address": ipAddress,
			"timestamp":  timestamp.Format(time.RFC3339),
		},
	)
}

func (l *Log) SystemStartup(timestamp time.Time, version, environment string) {
	l.DetailedLog(
		timestamp,
		INFO,
		"Authentication service started",
		SYSTEM,
		"system",
		"system",
		map[string]interface{}{
			"version":      version,
			"environment":  environment,
			"startup_time": timestamp.Format(time.RFC3339),
		},
	)
}

func (l *Log) ToJSON() (string, error) {
	jsonData, err := json.Marshal(l)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
