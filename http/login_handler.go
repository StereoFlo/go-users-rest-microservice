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
	privateKey := getKeyData("private_key.pem")
	publicKey := getKeyData("public_key.pem")
	jwt := jwt_token.NewJWT(privateKey, publicKey)
	acExpire := time.Now().Add(10 * time.Hour)
	rtExpire := time.Now().Add(20 * time.Hour)
	accessToken := getToken(jwt, acExpire, user)
	refreshToken := getToken(jwt, rtExpire, user)
	token = &entity.Token{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpire:  acExpire,
		RefreshTokenExpire: rtExpire,
		UserId:             user.ID,
	}
	handler.UserApp.SaveToken(token)
	user, err = handler.UserApp.GetUser(user.ID, 1)
	context.JSON(http.StatusOK, handler.responder.Success(token))
}

func getKeyData(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	return data
}

func getToken(jwt jwt_token.JWT, time time.Time, user *entity.User) string {
	accessToken, err := jwt.Get(time, user)
	if err != nil {
		log.Fatalln(err)
	}
	return accessToken
}
