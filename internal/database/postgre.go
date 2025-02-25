package database

import (
	"auth-service/internal/user"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v3"
	"log"
	"os"
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

func CreateConnection() *pgx.Conn {
	config, err := LoadConfig("/home/zaproktnost/GolandProjects/log-aggretator/configs/config.yaml")
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

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Println("Error connecting to database:", err)
		return nil
	}
	return conn
}

func CreateTable(conn *pgx.Conn) {
	query := `CREATE TABLE IF NOT EXISTS users (id serial primary key, username varchar(60) unique, password varchar(60))`

	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
}

func CreateUser(conn *pgx.Conn, u user.RegularUser) {
	query := `INSERT INTO users (username, password) values (@username, @password)`

	args := pgx.NamedArgs{
		"username": u.Username,
		"password": u.Password.Hashed,
	}
	_, err := conn.Exec(context.Background(), query, args)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
}

func SelectUsers(conn *pgx.Conn, users []string) {
	query := `SELECT * FROM users`

	info, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	defer info.Close()
}
