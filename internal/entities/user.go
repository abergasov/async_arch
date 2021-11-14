package entities

import (
	"github.com/golang-jwt/jwt"
)

const (
	UserCUDBrokerTopic = "userStream"
	UserCreatedEvent   = "UserCreated"
	UserUpdatedEvent   = "UserUpdated"
)

//type UserAccount1 struct {
//	ID       int64     `json:"-" db:"user_id"`
//	PublicID uuid.UUID `json:"public_id" db:"public_id"`
//	UserName string    `json:"user_name" db:"user_name"`
//	UserMail string    `json:"user_mail" db:"user_mail"`
//	UserRole string    `json:"user_role" db:"user_role"`
//	Version  int       `json:"version" db:"user_version"`
//	Active   bool      `json:"active" db:"active"`
//}

type UserJWT struct {
	UserID      string `json:"id"`
	UserVersion int64  `json:"v"`
	jwt.StandardClaims
}
