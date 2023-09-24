package entity

import (
	"fmt"
	"github.com/badoux/checkmail"
	"strings"
	"time"
)

const (
	Login = "login"
)

type User struct {
	ID        int       `gorm:"primary_key;auto_increment" json:"id"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Tokens    []Token   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Users []User

func (u User) ValidateUser(action string) map[string]string {
	var errorMessages = make(map[string]string)
	switch strings.ToLower(action) {
	case Login:
		if u.Password == "" {
			errorMessages["password_required"] = "password is required"
		}
		if u.Email == "" {
			errorMessages["email_required"] = "email is required"
		}
		err := checkmail.ValidateFormat(u.Email)
		if err != nil {
			errorMessages["invalid_email"] = "please provide a valid email"
		}
	default:
		errorMessages["password_required"] = fmt.Sprintf("unknown action %s", action)
	}
	return errorMessages
}
