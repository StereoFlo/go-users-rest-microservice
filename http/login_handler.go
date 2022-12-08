package http

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"user-app/application"
	"user-app/entity"
	"user-app/infrastructure"
)

type LoginHandler struct {
	UserApp   *application.UserApp
	responder *infrastructure.Responder
}

func NewLoginHandler(userApp *application.UserApp, responder *infrastructure.Responder) *LoginHandler {
	return &LoginHandler{
		UserApp:   userApp,
		responder: responder,
	}
}

func (handler *LoginHandler) Login(context *gin.Context) {
	var reqUser *entity.User
	var token *entity.Token
	err := context.ShouldBindJSON(&reqUser)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, handler.responder.Fail("Invalid json provided"))
		return
	}
	validateUser := reqUser.Validate("login")
	if len(validateUser) > 0 {
		context.JSON(http.StatusUnprocessableEntity, handler.responder.Fail(validateUser))
		return
	}
	dbUser, err := handler.UserApp.GetUserByEmail(reqUser.Email)
	if err != nil {
		context.JSON(http.StatusNotFound, handler.responder.Fail("user not found"))
		return
	}

	err = infrastructure.VerifyPassword(reqUser.Password, dbUser.Password)
	if err != nil {
		context.JSON(http.StatusNotFound, handler.responder.Fail("password is wrong"))
		return
	}
	jwt := infrastructure.NewToken()
	acExpire := time.Now().Add(10 * time.Hour)
	rtExpire := time.Now().Add(20 * time.Hour)
	accessToken := getToken(jwt, acExpire, dbUser)
	refreshToken := getToken(jwt, rtExpire, dbUser)
	t, _ := jwt.Validate(accessToken)
	token = &entity.Token{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpire:  acExpire,
		RefreshTokenExpire: rtExpire,
		UserId:             dbUser.ID,
		UUID:               t.Data.TokenId,
	}
	_, err = handler.UserApp.SaveToken(token)
	if err != nil {
		context.JSON(http.StatusInternalServerError, handler.responder.Fail(err))
		return
	}
	context.JSON(http.StatusOK, handler.responder.Success(token))
}

func getToken(jwt infrastructure.Token, time time.Time, user *entity.User) string {
	accessToken, err := jwt.Get(time, user.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return accessToken
}
