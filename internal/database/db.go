package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v3"
	"log"
	"os"
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
}

func LoadConfig(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	c := new(Config)

	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func CreateConnection(filePath string) {
	once.Do(func() {
		config, err := LoadConfig(filePath)
		if err != nil {
			fmt.Println(err)
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
			log.Println("Error connecting to database:", err)
		}

		dbPool, err = pgxpool.NewWithConfig(context.Background(), poolCfg)
		if err != nil {
			log.Println("Error connecting to database:", err)
		}
	})
}

func GetConnection() *pgxpool.Pool {
	if dbPool == nil {
		log.Println("Database is not connected")
		return nil
	}
	return dbPool
}
