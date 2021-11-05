package user

import (
	"async_arch/internal/entities"
	"async_arch/internal/storage/database"
)

type UserRepo interface {
	AddUser(account entities.UserAccount) error
}

type User struct {
	conn database.DBConnector
}

func InitUserRepo(conn database.DBConnector) *User {
	return &User{conn: conn}
}

func (u *User) AddUser(account entities.UserAccount) error {
	//u.conn.Client().NamedExec("INSERT INTO")
	return nil
}

func (u *User) UpdateUser(account entities.UserAccount) error {
	return nil
}

func (u *User) DeleteUser(account entities.UserAccount) error {
	return nil
}
