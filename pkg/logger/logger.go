package logger

import "time"

type Log struct {
	Id        int       `json:"id"`
	TimeStamp time.Time `json:"TimeStamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
}

func (l *Log) CreateLog(t time.Time, level, message, username, role string) {
	l.TimeStamp = t
	l.Level = level
	l.Message = message
	l.Username = username
	l.Role = role
}
