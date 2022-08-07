package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
	"user-app/application"
	"user-app/entity"
	"user-app/infrastructure"
)

type LoginHandler struct {
	UserApp   application.UserApp
	responder *infrastructure.Responder
}

func NewLoginHandler(userApp application.UserApp, responder *infrastructure.Responder) *LoginHandler {
	return &LoginHandler{
		UserApp:   userApp,
		responder: responder,
	}
}

func (handler *LoginHandler) Login(context *gin.Context) {
	var user *entity.User
	var token *entity.Token
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, handler.responder.Fail("Invalid json provided"))
		return
	}
	validateUser := user.Validate("login")
	if len(validateUser) > 0 {
		context.JSON(http.StatusUnprocessableEntity, validateUser)
		return
	}

	passwordRaw := user.Password
	user, err = handler.UserApp.GetUserByEmail(user.Email)
	if err != nil {
		context.JSON(http.StatusNotFound, handler.responder.Fail("user not found"))
		return
	}

	err = infrastructure.VerifyPassword(user.Password, passwordRaw)
	if err != nil {
		context.JSON(http.StatusNotFound, handler.responder.Fail("password is wrong"))
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
	handler.UserApp.SaveToken(token)
	user, err = handler.UserApp.GetUser(user.ID, 1)
	context.JSON(http.StatusOK, handler.responder.Success(token))
}
