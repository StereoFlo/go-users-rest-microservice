package infrastructure

import (
	"fmt"
	"github.com/badoux/checkmail"
	"strings"
	"user-app/entity"
)

const (
	Login = "login"
)

func ValidateUser(user *entity.User, action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case Login:
		if user.Password == "" {
			errorMessages["password_required"] = "password is required"
		}
		if user.Email == "" {
			errorMessages["email_required"] = "email is required"
		}
		err = checkmail.ValidateFormat(user.Email)
		if err != nil {
			errorMessages["invalid_email"] = "please provide a valid email"
		}
	default:
		errorMessages["password_required"] = fmt.Sprintf("unknown action %s", action)
	}
	return errorMessages
}
