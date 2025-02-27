package user_db

import (
	"auth-service/internal/database"
	"auth-service/internal/user"
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"strings"
)

func CreateTable() {
	conn := database.GetConnection()
	query := `CREATE TABLE IF NOT EXISTS users (id serial primary key, username varchar(60) unique, password varchar(60))`

	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
}

func CreateUser(u user.RegularUser) {
	conn := database.GetConnection()
	query := `INSERT INTO users (username, password) values (@username, @password)`

	args := pgx.NamedArgs{
		"username": strings.TrimSpace(u.Username),
		"password": u.Password.Hashed,
	}
	_, err := conn.Exec(context.Background(), query, args)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
}

func SelectUsers(users []string) {
	conn := database.GetConnection()
	query := `SELECT * FROM users`

	info, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	defer info.Close()
}

func GetUser(username string) bool {
	conn := database.GetConnection()
	query := `SELECT username FROM users WHERE username=@username`

	args := pgx.NamedArgs{
		"username": username,
	}

	row := conn.QueryRow(context.Background(), query, args)

	var u string

	err := row.Scan(&u)
	if err != nil {
		log.Println("Error")
	}

	return !(u == "")
}

func GetUserPassword(username string) string {
	conn := database.GetConnection()
	query := `SELECT password FROM users WHERE username=@username`

	args := pgx.NamedArgs{
		"username": username,
	}

	row := conn.QueryRow(context.Background(), query, args)

	var p string

	err := row.Scan(&p)
	if err != nil {
		log.Println("Error")
	}

	return p
}
