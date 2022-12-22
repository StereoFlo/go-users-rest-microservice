package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"user-app/application"
	"user-app/entity"
	"user-app/infrastructure"
)

type UserHandler struct {
	App       *application.UserApp
	responder *infrastructure.Responder
}

func NewUserHandler(userApp *application.UserApp, responder *infrastructure.Responder) *UserHandler {
	return &UserHandler{
		App:       userApp,
		responder: responder,
	}
}

func (handler *UserHandler) SaveUser(context *gin.Context) {
	var user entity.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, handler.responder.Fail(gin.H{
			"invalid_json": "invalid json",
		}))
		return
	}
	_, err = handler.App.SaveUser(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, handler.responder.Fail(err))
		return
	}
	context.JSON(http.StatusCreated, handler.responder.Success(user))
}

func (handler *UserHandler) GetList(context *gin.Context) {
	users := entity.Users{}
	limit, _ := strconv.Atoi(context.Query("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(context.Query("offset"))
	cnt, _ := handler.App.GetUserCount()
	if cnt == 0 {
		context.JSON(http.StatusOK, handler.responder.SuccessList(0, limit, offset, gin.H{}))
		return
	}
	var err error
	users, err = handler.App.GetList(limit, offset)
	if err != nil {
		context.JSON(http.StatusInternalServerError, handler.responder.Fail(err))
		return
	}
	context.JSON(http.StatusOK, handler.responder.SuccessList(cnt, limit, offset, users))
}

func (handler *UserHandler) GetUser(context *gin.Context) {
	userId, err := strconv.Atoi(context.Param("user_id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	user, err := handler.App.GetUser(userId, 1)
	if err != nil {
		context.JSON(http.StatusInternalServerError, handler.responder.Fail(err))
		return
	}
	context.JSON(http.StatusOK, handler.responder.Success(user))
}
