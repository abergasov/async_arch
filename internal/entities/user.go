package entities

import "github.com/google/uuid"

type UserAccount struct {
	ID       int64     `json:"-" db:"user_id"`
	PublicID uuid.UUID `json:"public_id" db:"public_id"`
	UserName string    `json:"user_name" db:"user_name"`
	UserMail string    `json:"user_mail" db:"user_mail"`
	UserRole string    `json:"user_role" db:"user_role"`
	Active   bool      `json:"active" db:"active"`
}
