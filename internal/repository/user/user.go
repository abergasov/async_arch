package user

import (
	"database/sql"
	"errors"

	"async_arch/internal/entities"
	"async_arch/internal/logger"
	"async_arch/internal/storage/database"

	"github.com/google/uuid"
)

type UserRepo interface {
	GetUserByMail(email string) (usr *entities.UserAccount, err error)
	AddUser(googleUser *entities.GoogleUser, userRole string) error
	GetUserByPublicID(publicID uuid.UUID, version int) (*entities.UserAccount, error)
	ChangeRole(publicID uuid.UUID, version int, role string) (*entities.UserAccount, error)
}

type User struct {
	conn database.DBConnector
}

func InitUserRepo(conn database.DBConnector) *User {
	return &User{conn: conn}
}

func (u *User) GetUserByMail(email string) (*entities.UserAccount, error) {
	sqlS := "SELECT user_id, public_id, user_mail, user_name, user_version, user_role, user_role, active FROM users WHERE user_mail = $1"
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
	sqlS := "SELECT user_id, public_id, user_mail, user_name, user_version, user_role, user_role, active FROM users WHERE public_id = $1 AND user_version = $2"
	var usr entities.UserAccount
	err := u.conn.Client().QueryRowx(sqlS, publicID, version).StructScan(&usr)
	return &usr, err
}

func (u *User) ChangeRole(publicID uuid.UUID, version int, role string) (*entities.UserAccount, error) {
	sqlU := "UPDATE users SET user_role = $1, user_version = user_version + 1 WHERE public_id = $2 AND user_version = $3"
	rows, err := u.conn.Client().Exec(sqlU, role, publicID, version)
	if err != nil {
		logger.Error("error change role", err)
		return nil, err
	}
	i, _ := rows.RowsAffected()
	if i == 0 {
		err = errors.New("user not updated")
		logger.Error("error change role", err)
		return nil, err
	}
	return u.GetUserByPublicID(publicID, version+1)
}

func (u *User) CreateUser(account *entities.UserAccount) error {
	_, err := u.conn.Client().NamedExec(
		"INSERT INTO users (public_id, user_mail, user_name, user_version, user_role, active) VALUES (:public_id, :user_mail, :user_name, :user_version, :user_role, :active)",
		map[string]interface{}{
			"public_id":    account.PublicID,
			"user_mail":    account.UserMail,
			"user_name":    account.UserName,
			"user_version": account.Version,
			"user_role":    account.UserRole,
			"active":       account.Active,
		})
	if err != nil {
		logger.Error("error insert user", err)
		return err
	}
	return nil
}

func (u *User) DeleteUser(account entities.UserAccount) error {
	return nil
}
