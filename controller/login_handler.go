package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
	"user-app/application"
	"user-app/entity"
	"user-app/infrastructure"
)

type Authenticate struct {
	UserApp application.UserApp
}

func NewAuth(userApp application.UserApp) *Authenticate {
	return &Authenticate{
		UserApp: userApp,
	}
}

func (authInterface *Authenticate) Login(context *gin.Context) {
	var user *entity.User
	var token *entity.Token
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
	user, err = authInterface.UserApp.GetUserByEmail(user.Email)
	if err != nil {
		context.JSON(http.StatusNotFound, "user not found")
		return
	}

	err = infrastructure.VerifyPassword(user.Password, passwordRaw)
	if err != nil {
		context.JSON(http.StatusNotFound, "password is wrong")
		return
	}
	now := time.Now()
	token = &entity.Token{
		AccessToken:        uuid.New().String(),
		RefreshToken:       uuid.New().String(),
		AccessTokenExpire:  now.Add(time.Hour * 8),
		RefreshTokenExpire: now.Add(time.Hour * 16),
		UserId:             user.ID,
	}
	authInterface.UserApp.SaveToken(token)
	user, err = authInterface.UserApp.GetUser(user.ID, 1)
	context.JSON(http.StatusOK, token)
}
