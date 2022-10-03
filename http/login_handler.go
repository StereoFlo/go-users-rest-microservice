package http

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
	"user-app/application"
	"user-app/entity"
	"user-app/infrastructure"
	jwt_token "user-app/infrastructure/jwt-token"
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
	user, err = handler.UserApp.GetUserByEmail(user.Email)
	if err != nil {
		context.JSON(http.StatusNotFound, handler.responder.Fail("user not found"))
		return
	}

	passwordRaw := user.Password
	err = infrastructure.VerifyPassword(user.Password, passwordRaw)
	if err != nil {
		context.JSON(http.StatusNotFound, handler.responder.Fail("password is wrong"))
		return
	}
	privateKey, err := os.ReadFile("private_key.pem")
	if err != nil {
		log.Fatalln(err)
	}
	publicKey, err := os.ReadFile("public_key.pem")
	if err != nil {
		log.Fatalln(err)
	}
	jwt := jwt_token.NewJWT(privateKey, publicKey)
	now := time.Now()
	accessToken, err := jwt.Get(now.Sub(now.Add(time.Hour*8)), user)
	if err != nil {
		log.Fatalln(err)
	}
	refreshToken, err := jwt.Get(now.Sub(now.Add(time.Hour*16)), user)
	if err != nil {
		log.Fatalln(err)
	}
	token = &entity.Token{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpire:  now.Add(time.Hour * 8),
		RefreshTokenExpire: now.Add(time.Hour * 16),
		UserId:             user.ID,
	}
	handler.UserApp.SaveToken(token)
	user, err = handler.UserApp.GetUser(user.ID, 1)
	context.JSON(http.StatusOK, handler.responder.Success(token))
}
