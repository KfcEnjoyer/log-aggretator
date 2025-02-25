package user

type User interface {
	LogIn() error
	GetRole(username string) string
}
type RegularUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Admin struct {
	RegularUser
	Permissions []string
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
