package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
)

var conn *pgx.Conn
var once sync.Once

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

func CreateConnection() {
	once.Do(func() {
		config, err := LoadConfig("C:\\Users\\user\\Desktop\\coding and stuff\\study or portfolio projects\\auth-service\\configs\\config.yaml")
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

		conn, err = pgx.Connect(context.Background(), connStr)
		if err != nil {
			log.Println("Error connecting to database:", err)
		}
	})
}

func GetConnection() *pgx.Conn {
	if conn == nil {
		CreateConnection()
	}
	return conn
}
