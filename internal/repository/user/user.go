package user

import (
	"database/sql"

	"async_arch/internal/entities"
	"async_arch/internal/logger"
	"async_arch/internal/storage/database"

	"github.com/google/uuid"
)

type UserRepo interface {
	GetUserByMail(email string) (usr *entities.UserAccount, err error)
	AddUser(googleUser *entities.GoogleUser, userRole string) error
	GetUserByPublicID(publicID uuid.UUID, version int) (*entities.UserAccount, error)
}

type User struct {
	conn database.DBConnector
}

func InitUserRepo(conn database.DBConnector) *User {
	return &User{conn: conn}
}

func (u *User) GetUserByMail(email string) (*entities.UserAccount, error) {
	sqlS := "SELECT user_id, public_id, user_mail, user_name, user_version, user_role, user_role FROM users WHERE user_mail = $1"
	var usr entities.UserAccount
	err := u.conn.Client().QueryRowx(sqlS, email).StructScan(&usr)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return &usr, err
}

func (u *User) AddUser(googleUser *entities.GoogleUser, userRole string) error {
	_, err := u.conn.Client().NamedExec(
		"INSERT INTO users (public_id, user_mail, user_name, user_version, user_role, active) VALUES (:public_id, :user_mail, :user_name, :user_version, :user_role, :active)",
		map[string]interface{}{
			"public_id":    uuid.New(),
			"user_mail":    googleUser.Email,
			"user_name":    googleUser.Name,
			"user_version": 1,
			"user_role":    userRole,
			"active":       1,
		})
	if err != nil {
		logger.Error("error insert user", err)
		return err
	}
	return nil
}

func (u *User) GetUserByPublicID(publicID uuid.UUID, version int) (*entities.UserAccount, error) {
	sqlS := "SELECT user_id, public_id, user_mail, user_name, user_version, user_role, user_role FROM users WHERE public_id = $1 AND user_version = $2"
	var usr entities.UserAccount
	err := u.conn.Client().QueryRowx(sqlS, publicID, version).StructScan(&usr)
	return &usr, err
}

func (u *User) UpdateUser(account entities.UserAccount) error {
	return nil
}

func (u *User) DeleteUser(account entities.UserAccount) error {
	return nil
}
