package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"user-app/application"
	"user-app/entity"
)

type Users struct {
	appInterface application.UserAppInterface
}

func NewUsers(userAppInterface application.UserAppInterface) *Users {
	return &Users{
		appInterface: userAppInterface,
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
	newUser, err := userStr.appInterface.SaveUser(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusCreated, newUser.PublicUser())
}

func (userStr *Users) GetList(context *gin.Context) {
	users := entity.Users{}
	var err error
	users, err = userStr.appInterface.GetList()
	if err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	context.JSON(http.StatusOK, users.PublicUsers())
}

func (userStr *Users) GetUser(context *gin.Context) {
	userId, err := strconv.ParseUint(context.Param("user_id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := userStr.appInterface.GetUser(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	context.JSON(http.StatusOK, user.PublicUser())
}
