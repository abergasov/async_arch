package entities

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const (
	UserCUDBrokerTopic = "userStream"
	UserCreatedEvent   = "UserCreated"
	UserUpdatedEvent   = "UserUpdated"
)

type UserAccount struct {
	ID       int64     `json:"-" db:"user_id"`
	PublicID uuid.UUID `json:"public_id" db:"public_id"`
	UserName string    `json:"user_name" db:"user_name"`
	UserMail string    `json:"user_mail" db:"user_mail"`
	UserRole string    `json:"user_role" db:"user_role"`
	Version  int       `json:"version" db:"user_version"`
	Active   bool      `json:"active" db:"active"`
}

type UserJWT struct {
	UserID      uuid.UUID `json:"id"`
	UserVersion int       `json:"v"`
	jwt.StandardClaims
}
