package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"user-app/internal/application"
	entity2 "user-app/internal/entity"
	"user-app/pkg/utils"
)

type LoginHandler struct {
	userApp   *application.UserApp
	responder *utils.Responder
	token     utils.Token
}

func NewLoginHandler(userApp *application.UserApp, responder *utils.Responder, token utils.Token) *LoginHandler {
	return &LoginHandler{
		userApp:   userApp,
		responder: responder,
		token:     token,
	}
}

func (handler *LoginHandler) Login(ctx *gin.Context) {
	var reqUser *entity2.User
	err := ctx.ShouldBindJSON(&reqUser)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, handler.responder.Fail("Invalid json provided"))
		return
	}
	validateUser := reqUser.ValidateUser(entity2.Login)
	if len(validateUser) > 0 {
		ctx.JSON(http.StatusUnprocessableEntity, handler.responder.Fail(validateUser))
		return
	}
	err, dbUser := handler.userApp.GetUserByEmail(reqUser.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, handler.responder.Fail("user not found"))
		return
	}

	err = utils.VerifyPassword(dbUser.Password, reqUser.Password)
	if err != nil {
		ctx.JSON(http.StatusNotFound, handler.responder.Fail("password is wrong"))
		return
	}
	token, err := handler.makeNewToken(dbUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, handler.responder.Fail(err))
		return
	}

	ctx.JSON(http.StatusOK, handler.responder.Success(token))
}

func (handler *LoginHandler) makeNewToken(dbUser *entity2.User) (*entity2.Token, error) {
	var token *entity2.Token
	acExpire := time.Now().Add(10 * time.Hour)
	rtExpire := time.Now().Add(20 * time.Hour)
	accessToken, err := getToken(handler.token, acExpire, dbUser)
	if err != nil {
		return nil, err
	}

	refreshToken, err := getToken(handler.token, rtExpire, dbUser)
	if err != nil {
		return nil, err
	}

	t, err := handler.token.Validate(*accessToken)
	if err != nil {
		return nil, err
	}

	token = &entity2.Token{
		AccessToken:        *accessToken,
		RefreshToken:       *refreshToken,
		AccessTokenExpire:  acExpire,
		RefreshTokenExpire: rtExpire,
		UserId:             dbUser.ID,
		UUID:               t.Data.TokenId,
	}
	err, _ = handler.userApp.SaveToken(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func getToken(jwt utils.Token, time time.Time, user *entity2.User) (*string, error) {
	token, err := jwt.Get(time, user.ID)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
