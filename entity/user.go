package entity

import (
	"github.com/badoux/checkmail"
	"strings"
	"time"
	"user-app/infrastructure/security"
)

type User struct {
	ID           int       `gorm:"primary_key;auto_increment" json:"id"`
	Email        string    `gorm:"size:100;not null;unique" json:"email"`
	Password     string    `gorm:"size:100;not null;" json:"password"`
	AccessTokens []Token   `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type PublicUser struct {
	ID           int     `gorm:"primary_key;auto_increment" json:"id"`
	Email        string  `gorm:"size:100;not null;" json:"email"`
	AccessTokens []Token `json:"access_tokens"`
}

func (user *User) BeforeSave() error {
	hashPassword, err := security.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashPassword)
	return nil
}

type Users []User

func (users Users) PublicUsers() []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.PublicUser()
	}
	return result
}

func (user *User) PublicUser() interface{} {
	return &PublicUser{
		ID:           user.ID,
		Email:        user.Email,
		AccessTokens: user.AccessTokens,
	}
}
func (u *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		if u.Email == "" {
			errorMessages["email_required"] = "email required"
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["invalid_email"] = "email email"
			}
		}

	case "login":
		if u.Password == "" {
			errorMessages["password_required"] = "password is required"
		}
		if u.Email == "" {
			errorMessages["email_required"] = "email is required"
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	case "forgot_password":
		if u.Email == "" {
			errorMessages["email_required"] = "email required"
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	default:
		if u.Password == "" {
			errorMessages["password_required"] = "password is required"
		}
		if u.Password != "" && len(u.Password) < 6 {
			errorMessages["invalid_password"] = "password should be at least 6 characters"
		}
		if u.Email == "" {
			errorMessages["email_required"] = "email is required"
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	}
	return errorMessages
}
