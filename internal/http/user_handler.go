package http

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"user-app/internal/application"
	"user-app/internal/entity"
	"user-app/pkg/utils"
)

type UserHandler struct {
	userApp   *application.UserApp
	responder *utils.Responder
}

func NewUserHandler(userApp *application.UserApp, responder *utils.Responder) *UserHandler {
	return &UserHandler{
		userApp:   userApp,
		responder: responder,
	}
}

func (handler *UserHandler) SaveUser(ctx *gin.Context) {
	var user entity.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, handler.responder.Fail("invalid json"))
		return
	}
	bytePass, err := utils.Hash(user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, handler.responder.Fail(err))
		return
	}
	user.Password = string(bytePass)
	err, _ = handler.userApp.SaveUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, handler.responder.Fail(err))
		return
	}
	ctx.JSON(http.StatusCreated, handler.responder.Success(user))
}

func (handler *UserHandler) GetList(ctx *gin.Context) {
	users := entity.Users{}
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, handler.responder.Fail(err))
		return
	}

	if limit == 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, handler.responder.Fail(err))
		return
	}

	err, cnt := handler.userApp.GetUserCount()
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, handler.responder.Fail(err))
		return
	}

	if *cnt == 0 {
		log.Print(err)
		ctx.JSON(http.StatusOK, handler.responder.SuccessList(0, limit, offset, gin.H{}))
		return
	}

	err, users = handler.userApp.GetList(limit, offset)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, handler.responder.Fail(err))
		return
	}

	ctx.JSON(http.StatusOK, handler.responder.SuccessList(*cnt, limit, offset, users))
}

func (handler *UserHandler) GetUser(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	err, user := handler.userApp.GetUser(userId, 1)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, handler.responder.Fail(err))
		return
	}

	ctx.JSON(http.StatusOK, handler.responder.Success(user))
}
