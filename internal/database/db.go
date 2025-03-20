package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	dbPool *pgxpool.Pool
	once   sync.Once
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	} `yaml:"database"`
	Kafka struct {
		Brokers []string `yaml:"brokers"`
		Topic   string   `yaml:"topic"`
		GroupID string   `yaml:"groupID"`
	} `yaml:"kafka"`
}

func LoadConfig(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	c := new(Config)

	dataStr := string(data)
	dataStr = os.ExpandEnv(dataStr)

	err = yaml.Unmarshal([]byte(dataStr), c)
	if err != nil {
		return nil, err
	}

	if len(c.Kafka.Brokers) == 1 && strings.Contains(c.Kafka.Brokers[0], ",") {
		c.Kafka.Brokers = strings.Split(c.Kafka.Brokers[0], ",")
	}

	return c, nil
}

func CreateConnection(filePath string) {
	once.Do(func() {
		config, err := LoadConfig(filePath)
		if err != nil {
			log.Printf("Error loading config: %v", err)
			return
		}

		connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.Dbname,
		)

		poolCfg, err := pgxpool.ParseConfig(connStr)
		if err != nil {
			log.Println("Error parsing pool config:", err)
			return
		}

		dbPool, err = pgxpool.NewWithConfig(context.Background(), poolCfg)
		if err != nil {
			log.Println("Error connecting to database:", err)
			return
		}

		if err := dbPool.Ping(context.Background()); err != nil {
			log.Println("Error pinging database:", err)
			return
		}

		log.Println("Successfully connected to database")
	})
}

func GetConnection() *pgxpool.Pool {
	if dbPool == nil {
		log.Println("Database is not connected")
		return nil
	}
	return dbPool
}
