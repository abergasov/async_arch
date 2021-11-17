package user

import (
	"database/sql"
	"errors"

	"async_arch/internal/entities"
	"async_arch/internal/logger"
	"async_arch/internal/storage/database"

	"github.com/abergasov/schema_registry/pkg/grpc/user"
	"github.com/google/uuid"
)

type UserRepo interface {
	GetUserByMail(email string) (usr *user.UserAccountV1, err error)
	AddUser(googleUser *entities.GoogleUser, userRole string) error
	GetByPublicID(publicID string, version int64) (*user.UserAccountV1, error)
	ChangeRole(publicID string, version int64, role string) (*user.UserAccountV1, error)
}

type User struct {
	conn database.DBConnector
}

func InitUserRepo(conn database.DBConnector) *User {
	return &User{conn: conn}
}

func (u *User) GetUserByMail(email string) (*user.UserAccountV1, error) {
	sqlS := "SELECT public_id, user_mail, user_name, user_version, user_role, active FROM users WHERE user_mail = $1"
	var usr user.UserAccountV1
	err := u.conn.Client().QueryRowx(sqlS, email).
		Scan(&usr.PublicID, &usr.UserMail, &usr.UserName, &usr.Version, &usr.UserRole, &usr.Active)
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

func (u *User) GetByPublicID(publicID string, version int64) (*user.UserAccountV1, error) {
	sqlS := "SELECT public_id, user_mail, user_name, user_version, user_role, active FROM users WHERE public_id = $1 AND user_version = $2"
	var usr user.UserAccountV1
	err := u.conn.Client().QueryRowx(sqlS, publicID, version).
		Scan(&usr.PublicID, &usr.UserMail, &usr.UserName, &usr.Version, &usr.UserRole, &usr.Active)
	return &usr, err
}

func (u *User) ChangeRole(publicID string, version int64, role string) (*user.UserAccountV1, error) {
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
	return u.GetByPublicID(publicID, version+1)
}

func (u *User) CreateUser(account *user.UserAccountV1) error {
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

func (u *User) UpdateUser(account *user.UserAccountV1) error {
	if _, err := u.conn.Client().NamedExec(
		"UPDATE users SET user_mail=:user_mail, user_name=:user_name, user_version=:user_version, user_role=:user_role, active=:active WHERE public_id=:public_id",
		map[string]interface{}{
			"public_id":    account.PublicID,
			"user_mail":    account.UserMail,
			"user_name":    account.UserName,
			"user_version": account.Version,
			"user_role":    account.UserRole,
			"active":       account.Active,
		}); err != nil {
		logger.Error("error insert user", err)
		return err
	}
	return nil
}

func (u *User) GetActiveWorkers() ([]uuid.UUID, error) {
	result := make([]uuid.UUID, 0, 1000)
	rows, err := u.conn.Client().Queryx("SELECT public_id FROM users WHERE active = true AND user_role = 'worker'")
	if err != nil {
		logger.Error("error get active workers", err)
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var uID uuid.UUID
		rows.Scan(&uID)
		result = append(result, uID)
	}
	return result, nil
}
