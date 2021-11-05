package user

import "async_arch/internal/storage/database"

type UserRepo interface {
	AddUser()
}

type User struct {
	conn database.DBConnector
}

func InitUserRepo(conn database.DBConnector) *User {
	return &User{conn: conn}
}

func (u *User) AddUser() {}
