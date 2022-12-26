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
	userApp   *application.UserApp
	responder *infrastructure.Responder
	jwtToken  *infrastructure.Token
}

func NewLoginHandler(userApp *application.UserApp, responder *infrastructure.Responder, jwtToken *infrastructure.Token) *LoginHandler {
	return &LoginHandler{
		userApp:   userApp,
		responder: responder,
		jwtToken:  jwtToken,
	}
}

func (handler *LoginHandler) Login(context *gin.Context) {
	var reqUser *entity.User
	err := context.ShouldBindJSON(&reqUser)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, handler.responder.Fail("Invalid json provided"))
		return
	}
	validateUser := infrastructure.ValidateUser(reqUser, infrastructure.Login)
	if len(validateUser) > 0 {
		context.JSON(http.StatusUnprocessableEntity, handler.responder.Fail(validateUser))
		return
	}
	dbUser, err := handler.userApp.GetUserByEmail(reqUser.Email)
	if err != nil {
		context.JSON(http.StatusNotFound, handler.responder.Fail("user not found"))
		return
	}

	err = infrastructure.VerifyPassword(dbUser.Password, reqUser.Password)
	if err != nil {
		context.JSON(http.StatusNotFound, handler.responder.Fail("password is wrong"))
		return
	}
	token, done := handler.makeNewToken(context, dbUser)
	if done {
		return
	}
	context.JSON(http.StatusOK, handler.responder.Success(token))
}

func (handler *LoginHandler) makeNewToken(context *gin.Context, dbUser *entity.User) (*entity.Token, bool) {
	var token *entity.Token
	acExpire := time.Now().Add(10 * time.Hour)
	rtExpire := time.Now().Add(20 * time.Hour)
	accessToken := getToken(handler.jwtToken, acExpire, dbUser)
	refreshToken := getToken(handler.jwtToken, rtExpire, dbUser)
	t, _ := handler.jwtToken.Validate(accessToken)
	token = &entity.Token{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpire:  acExpire,
		RefreshTokenExpire: rtExpire,
		UserId:             dbUser.ID,
		UUID:               t.Data.TokenId,
	}
	_, err := handler.userApp.SaveToken(token)
	if err != nil {
		context.JSON(http.StatusInternalServerError, handler.responder.Fail(err))
		return nil, true
	}
	return token, false
}

func getToken(jwt *infrastructure.Token, time time.Time, user *entity.User) string {
	accessToken, err := jwt.Get(time, user.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return accessToken
}
