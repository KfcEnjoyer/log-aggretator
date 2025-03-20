package database

import (
	"auth-service/internal/user"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strings"
)

type UserRepo interface {
	CreateUser(u user.RegularUser) error
	GetUserByUsername(username string) (user.RegularUser, error)
}

type UserRepoService struct {
	dbPool *pgxpool.Pool
}

func NewUserRepoService() *UserRepoService {
	return &UserRepoService{
		dbPool: GetConnection(),
	}
}

func (r *UserRepoService) CreateUser(u user.RegularUser) error {
	query := `INSERT INTO users (username, password, role) values (@username, @password, @role)`

	args := pgx.NamedArgs{
		"username": strings.TrimSpace(u.Username),
		"password": u.Password.Hashed,
		"role":     u.Role,
	}
	_, err := r.dbPool.Exec(context.Background(), query, args)
	if err != nil {
		return err
	}

	return nil
}

//func CreateTable() {
//	conn := database.GetConnection()
//	query := `CREATE TABLE IF NOT EXISTS users (id serial primary key, username varchar(60) unique, password varchar(60), role varchar(20) default "regular"`
//
//	_, err := conn.Exec(context.Background(), query)
//	if err != nil {
//		log.Fatalf("Query failed: %v", err)
//	}
//}

//func SelectUsers(users []string) {
//	conn := database.GetConnection()
//	query := `SELECT * FROM users`
//
//	info, err := conn.Query(context.Background(), query)
//	if err != nil {
//		log.Fatalf("Query failed: %v", err)
//	}
//
//	defer info.Close()
//}

func (r *UserRepoService) GetUserByUsername(username string) (user.RegularUser, error) {
	query := `SELECT * FROM users WHERE username=@username`

	args := pgx.NamedArgs{
		"username": username,
	}

	row := r.dbPool.QueryRow(context.Background(), query, args)

	var u user.RegularUser

	err := row.Scan(&u.Id, &u.Username, &u.Password.Hashed, &u.Role)
	if err != nil {
		log.Println("Error")
		return user.RegularUser{}, err
	}

	return u, nil
}
