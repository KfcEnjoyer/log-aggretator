package database

import (
	"auth-service/pkg/logger"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type LogRepo interface {
	CreateLog(l logger.Log) error
	GetRecentLogs(limit int) ([]logger.Log, error)
	GetLogsByLevel(level string, limit int) ([]logger.Log, error)
}

type LogRepoService struct {
	dbPool *pgxpool.Pool
}

func NewLogRepoService() *LogRepoService {
	return &LogRepoService{
		dbPool: GetConnection(),
	}
}

func (r *LogRepoService) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS logs (
		id SERIAL PRIMARY KEY, 
		level VARCHAR(20) NOT NULL, 
		message VARCHAR(255) NOT NULL, 
		category VARCHAR(50) NOT NULL,
		username VARCHAR(60) NOT NULL, 
		role VARCHAR(30) NOT NULL,
		metadata JSONB,
		timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`

	_, err := r.dbPool.Exec(context.Background(), query)
	if err != nil {
		log.Printf("Failed to create logs table: %v", err)
		return err
	}

	return nil
}

func (r *LogRepoService) CreateLog(l logger.Log) error {
	metadataJSON, err := json.Marshal(l.Metadata)
	if err != nil {
		log.Printf("Failed to convert metadata to JSON: %v", err)
		return err
	}

	query := `
	INSERT INTO logs (level, message, category, username, role, metadata, timestamp) 
	VALUES (@level, @message, @category, @username, @role, @metadata, @timestamp)
	`

	args := pgx.NamedArgs{
		"level":     l.Level,
		"message":   l.Message,
		"category":  l.Category,
		"username":  l.Username,
		"role":      l.Role,
		"metadata":  string(metadataJSON),
		"timestamp": l.TimeStamp,
	}

	_, err = r.dbPool.Exec(context.Background(), query, args)
	if err != nil {
		log.Printf("Query failed: %v\n", err)
		return err
	}

	return nil
}

func (r *LogRepoService) GetRecentLogs(limit int) ([]logger.Log, error) {
	query := `
	SELECT level, message, category, username, role, metadata, timestamp
	FROM logs
	ORDER BY timestamp DESC
	LIMIT $1
	`

	rows, err := r.dbPool.Query(context.Background(), query, limit)
	if err != nil {
		log.Printf("Query failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var logs []logger.Log
	for rows.Next() {
		var l logger.Log
		var metadataJSON string

		err = rows.Scan(&l.Level, &l.Message, &l.Category, &l.Username, &l.Role, &metadataJSON, &l.TimeStamp)
		if err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			continue
		}

		if metadataJSON != "" {
			if err = json.Unmarshal([]byte(metadataJSON), &l.Metadata); err != nil {
				log.Printf("Failed to parse metadata JSON: %v\n", err)
				l.Metadata = make(map[string]interface{})
			}
		}

		logs = append(logs, l)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v\n", err)
		return logs, err
	}

	return logs, nil
}

func (r *LogRepoService) GetLogsByLevel(level string, limit int) ([]logger.Log, error) {
	query := `
	SELECT level, message, category, username, role, metadata, timestamp
	FROM logs
	WHERE level = $1
	ORDER BY timestamp DESC
	LIMIT $2
	`

	rows, err := r.dbPool.Query(context.Background(), query, level, limit)
	if err != nil {
		log.Printf("Query failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var logs []logger.Log
	for rows.Next() {
		var l logger.Log
		var metadataJSON string

		err = rows.Scan(&l.Level, &l.Message, &l.Category, &l.Username, &l.Role, &metadataJSON, &l.TimeStamp)
		if err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			continue
		}

		if metadataJSON != "" {
			if err = json.Unmarshal([]byte(metadataJSON), &l.Metadata); err != nil {
				log.Printf("Failed to parse metadata JSON: %v\n", err)
				l.Metadata = make(map[string]interface{})
			}
		}

		logs = append(logs, l)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v\n", err)
		return logs, err
	}

	return logs, nil
}

func (r *LogRepoService) GetLogsByUsername(username string, limit int) ([]logger.Log, error) {
	query := `
	SELECT level, message, category, username, role, metadata, timestamp
	FROM logs
	WHERE username = $1
	ORDER BY timestamp DESC
	LIMIT $2
	`

	rows, err := r.dbPool.Query(context.Background(), query, username, limit)
	if err != nil {
		log.Printf("Query failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var logs []logger.Log
	for rows.Next() {
		var l logger.Log
		var metadataJSON string

		err = rows.Scan(&l.Level, &l.Message, &l.Category, &l.Username, &l.Role, &metadataJSON, &l.TimeStamp)
		if err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			continue
		}

		if metadataJSON != "" {
			if err = json.Unmarshal([]byte(metadataJSON), &l.Metadata); err != nil {
				log.Printf("Failed to parse metadata JSON: %v\n", err)
				l.Metadata = make(map[string]interface{})
			}
		}

		logs = append(logs, l)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v\n", err)
		return logs, err
	}

	return logs, nil
}
