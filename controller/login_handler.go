package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-app/application"
)

type Authenticate struct {
	userAppInterface application.UserAppInterface
}

func NewAuthenticate(userAppInterface application.UserAppInterface) *Authenticate {
	return &Authenticate{
		userAppInterface: userAppInterface,
	}
}

func (authInterface *Authenticate) Login(context *gin.Context) {
	userData := make(map[string]interface{})
	context.JSON(http.StatusOK, userData)
}

func (authInterface *Authenticate) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, "Successfully logged out")
}
