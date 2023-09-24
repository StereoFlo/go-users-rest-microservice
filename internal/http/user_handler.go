package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"user-app/internal/application"
	"user-app/internal/entity"
	"user-app/internal/infrastructure"
)

type UserHandler struct {
	userApp   *application.UserApp
	responder *infrastructure.Responder
}

func NewUserHandler(userApp *application.UserApp, responder *infrastructure.Responder) *UserHandler {
	return &UserHandler{
		userApp:   userApp,
		responder: responder,
	}
}

func (handler *UserHandler) SaveUser(ctx *gin.Context) {
	var user entity.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, handler.responder.Fail(gin.H{
			"invalid_json": "invalid json",
		}))
		return
	}
	_, err = handler.userApp.SaveUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, handler.responder.Fail(err))
		return
	}
	ctx.JSON(http.StatusCreated, handler.responder.Success(user))
}

func (handler *UserHandler) GetList(ctx *gin.Context) {
	users := entity.Users{}
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(ctx.Query("offset"))
	cnt, _ := handler.userApp.GetUserCount()
	if cnt == 0 {
		ctx.JSON(http.StatusOK, handler.responder.SuccessList(0, limit, offset, gin.H{}))
		return
	}
	var err error
	users, err = handler.userApp.GetList(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, handler.responder.Fail(err))
		return
	}
	ctx.JSON(http.StatusOK, handler.responder.SuccessList(cnt, limit, offset, users))
}

func (handler *UserHandler) GetUser(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	user, err := handler.userApp.GetUser(userId, 1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, handler.responder.Fail(err))
		return
	}
	ctx.JSON(http.StatusOK, handler.responder.Success(user))
}
