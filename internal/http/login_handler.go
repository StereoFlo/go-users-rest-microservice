package http

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"user-app/internal/application"
	entity2 "user-app/internal/entity"
	infrastructure2 "user-app/internal/infrastructure"
	"user-app/pkg/utils"
)

type LoginHandler struct {
	userApp   *application.UserApp
	responder *utils.Responder
}

func NewLoginHandler(userApp *application.UserApp, responder *utils.Responder) *LoginHandler {
	return &LoginHandler{
		userApp:   userApp,
		responder: responder,
	}
}

func (handler *LoginHandler) Login(context *gin.Context) {
	var reqUser *entity2.User
	err := context.ShouldBindJSON(&reqUser)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, handler.responder.Fail("Invalid json provided"))
		return
	}
	validateUser := infrastructure2.ValidateUser(reqUser, infrastructure2.Login)
	if len(validateUser) > 0 {
		context.JSON(http.StatusUnprocessableEntity, handler.responder.Fail(validateUser))
		return
	}
	err, dbUser := handler.userApp.GetUserByEmail(reqUser.Email)
	if err != nil {
		context.JSON(http.StatusNotFound, handler.responder.Fail("user not found"))
		return
	}

	err = utils.VerifyPassword(dbUser.Password, reqUser.Password)
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

func (handler *LoginHandler) makeNewToken(context *gin.Context, dbUser *entity2.User) (*entity2.Token, bool) {
	var token *entity2.Token
	jwt := utils.NewToken()
	acExpire := time.Now().Add(10 * time.Hour)
	rtExpire := time.Now().Add(20 * time.Hour)
	accessToken := getToken(jwt, acExpire, dbUser)
	refreshToken := getToken(jwt, rtExpire, dbUser)
	t, _ := jwt.Validate(accessToken)
	token = &entity2.Token{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpire:  acExpire,
		RefreshTokenExpire: rtExpire,
		UserId:             dbUser.ID,
		UUID:               t.Data.TokenId,
	}
	err, _ := handler.userApp.SaveToken(token)
	if err != nil {
		context.JSON(http.StatusInternalServerError, handler.responder.Fail(err))
		return nil, true
	}
	return token, false
}

func getToken(jwt utils.Token, time time.Time, user *entity2.User) string {
	accessToken, err := jwt.Get(time, user.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return accessToken
}
