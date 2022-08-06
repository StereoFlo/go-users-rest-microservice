package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"user-app/application"
	"user-app/entity"
	"user-app/infrastructure"
)

type Users struct {
	App       application.UserApp
	responder *infrastructure.Responder
}

func NewUsers(userApp application.UserApp, responder *infrastructure.Responder) *Users {
	return &Users{
		App:       userApp,
		responder: responder,
	}
}

func (userStr *Users) SaveUser(context *gin.Context) {
	var user entity.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "invalid json",
		})
		return
	}
	newUser, err := userStr.App.SaveUser(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, userStr.responder.Fail(err))
		return
	}
	context.JSON(http.StatusCreated, userStr.responder.Success(newUser.PublicUser()))
}

func (userStr *Users) GetList(context *gin.Context) {
	users := entity.Users{}
	limit, _ := strconv.Atoi(context.Query("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(context.Query("offset"))

	var err error
	users, err = userStr.App.GetList(limit, offset)
	if err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	context.JSON(http.StatusOK, userStr.responder.SuccessList(len(users.PublicUsers()), limit, offset, users))
}

func (userStr *Users) GetUser(context *gin.Context) {
	userId, err := strconv.Atoi(context.Param("user_id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := userStr.App.GetUser(userId, 1)
	if err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	context.JSON(http.StatusOK, userStr.responder.Success(user.PublicUser()))
}
