package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
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

func (lh *LoginHandler) Login(ctx *gin.Context) {
	var reqUser *entity2.User
	err := ctx.ShouldBindJSON(&reqUser)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusUnprocessableEntity, lh.responder.Fail("Invalid json provided"))
		return
	}
	validateUser := reqUser.ValidateUser(entity2.Login)
	if len(validateUser) > 0 {
		ctx.JSON(http.StatusUnprocessableEntity, lh.responder.Fail(validateUser))
		return
	}
	err, dbUser := lh.userApp.GetUserByEmail(ctx, reqUser.Email)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusNotFound, lh.responder.Fail("user not found"))
		return
	}

	err = utils.VerifyPassword(dbUser.Password, reqUser.Password)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusNotFound, lh.responder.Fail("password is wrong"))
		return
	}
	token, err := lh.makeNewToken(ctx, dbUser)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, lh.responder.Fail(err))
		return
	}

	ctx.JSON(http.StatusOK, lh.responder.Success(token))
}

func (lh *LoginHandler) makeNewToken(ctx context.Context, dbUser *entity2.User) (*entity2.Token, error) {
	var token *entity2.Token
	acExpire := time.Now().Add(10 * time.Hour)
	rtExpire := time.Now().Add(20 * time.Hour)
	accessToken, err := lh.token.Get(acExpire, dbUser.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := lh.token.Get(acExpire, dbUser.ID)
	if err != nil {
		return nil, err
	}

	t, err := lh.token.Validate(accessToken)
	if err != nil {
		return nil, err
	}

	token = &entity2.Token{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpire:  acExpire,
		RefreshTokenExpire: rtExpire,
		UserId:             dbUser.ID,
		UUID:               t.Data.TokenId,
	}
	err, _ = lh.userApp.SaveToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}
