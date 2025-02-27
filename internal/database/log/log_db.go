package log

import (
	"auth-service/internal/database"
	"auth-service/pkg/logger"
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

func CreateTable() {
	conn := database.GetConnection()
	query := `CREATE TABLE IF NOT EXISTS logs (id serial primary key, level varchar(20), message varchar(120), username varchar(60) unique, timestamp timestamptz NOT NULL default NOW())`

	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
}

func CreateLog(l logger.Log) {
	conn := database.GetConnection()
	query := `INSERT INTO logs (level, message, username, timestamp) values (@level, @message, @username, @timestamp)`

	args := pgx.NamedArgs{
		"level":     l.Level,
		"message":   l.Message,
		"username":  l.Username,
		"timestamp": l.TimeStamp,
	}

	_, err := conn.Exec(context.Background(), query, args)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
}
