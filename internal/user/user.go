package user

import (
	"auth-service/pkg/validator"
	"golang.org/x/crypto/bcrypt"
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
	h, err := bcrypt.GenerateFromPassword([]byte(p.Plain), bcrypt.DefaultCost)
	v.CheckError(err)

	p.Hashed = string(h)
}

//
//func GrantPermissions(user RegularUser) Admin{
//	var a a
//
//	if user.Role != "admin"{
//
//		return nil
//	}
//
//	a := Admin{
//		RegularUser: user,
//		Permissions: nil,
//	}
//
//	return a
//}
