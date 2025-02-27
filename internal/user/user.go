package user

import (
	"auth-service/pkg/validator"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

var v = validator.NewValidator()

type User interface {
	LogIn() error
	GetRole(username string) string
}
type RegularUser struct {
	Id       int      `json:"id"`
	Username string   `json:"username"`
	Password Password `json:"password"`
	Role     string   `json:"role"`
}

type Admin struct {
	RegularUser
	Permissions []string
}

type Password struct {
	Plain  string `json:"plain"`
	Hashed string
}

func (p *Password) Hash() {
	h, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(p.Plain)), bcrypt.DefaultCost)
	v.CheckError(err)

	p.Hashed = string(h)
}

func (p *Password) Compare(hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(p.Plain))
	if err != nil {
		log.Println(err)
	}

	return err == nil
}

//
//func GrantPermissions(user_db RegularUser) Admin{
//	var a a
//
//	if user_db.Role != "admin"{
//
//		return nil
//	}
//
//	a := Admin{
//		RegularUser: user_db,
//		Permissions: nil,
//	}
//
//	return a
//}
