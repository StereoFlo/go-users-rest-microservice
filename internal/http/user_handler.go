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

func (uh *UserHandler) SaveUser(ctx *gin.Context) {
	var user entity.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, uh.responder.Fail("invalid json"))
		return
	}
	bytePass, err := utils.Hash(user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, uh.responder.Fail(err))
		return
	}
	user.Password = string(bytePass)
	err, _ = uh.userApp.SaveUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, uh.responder.Fail(err))
		return
	}
	ctx.JSON(http.StatusCreated, uh.responder.Success(user))
}

func (uh *UserHandler) GetList(ctx *gin.Context) {
	users := entity.Users{}
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, uh.responder.Fail(err))
		return
	}

	if limit == 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, uh.responder.Fail(err))
		return
	}

	err, cnt := uh.userApp.GetUserCount()
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, uh.responder.Fail(err))
		return
	}

	if *cnt == 0 {
		log.Print(err)
		ctx.JSON(http.StatusOK, uh.responder.SuccessList(0, limit, offset, gin.H{}))
		return
	}

	err, users = uh.userApp.GetList(limit, offset)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, uh.responder.Fail(err))
		return
	}

	var count int
	count = int(*cnt)

	ctx.JSON(http.StatusOK, uh.responder.SuccessList(count, limit, offset, users))
}

func (uh *UserHandler) GetUser(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	err, user := uh.userApp.GetUser(userId, 1)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, uh.responder.Fail(err))
		return
	}

	ctx.JSON(http.StatusOK, uh.responder.Success(user))
}
