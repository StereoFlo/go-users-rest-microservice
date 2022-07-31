package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"user-app/application"
	"user-app/entity"
	"user-app/infrastructure/security"
)

type Authenticate struct {
	userAppInterface application.UserAppInterface
}

func NewAuthenticate(userAppInterface application.UserAppInterface) *Authenticate {
	return &Authenticate{
		userAppInterface: userAppInterface,
	}
}

func (authInterface *Authenticate) Login(context *gin.Context) {
	var user *entity.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	validateUser := user.Validate("login")
	if len(validateUser) > 0 {
		context.JSON(http.StatusUnprocessableEntity, validateUser)
		return
	}

	passwordRaw := user.Password
	user, err = authInterface.userAppInterface.GetUserByEmail(user.Email)
	if err != nil {
		context.JSON(http.StatusNotFound, "user not found")
		return
	}

	err = security.VerifyPassword(user.Password, passwordRaw)
	if err != nil {
		context.JSON(http.StatusNotFound, "password is wrong")
		return
	}
	token := entity.Token{AccessToken: "wwww", RefreshToken: "dddddd", AccessTokenExpire: time.Time{}, RefreshTokenExpire: time.Time{}, UserId: user.ID}
	user.AccessTokens = append(user.AccessTokens, token)
	fmt.Println(fmt.Sprintf("%#v", user))
	authInterface.userAppInterface.SaveUser(user)
	context.JSON(http.StatusOK, user)
}

func (authInterface *Authenticate) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, "Successfully logged out")
}
