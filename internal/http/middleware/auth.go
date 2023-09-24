package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"user-app/internal/application"
	"user-app/pkg/utils"
)

type Auth struct {
	userApp   *application.UserApp
	responder *utils.Responder
	token     utils.Token
}

func NewAuth(userApp *application.UserApp, responder *utils.Responder, token utils.Token) *Auth {
	return &Auth{userApp, responder, token}
}

func (userApp *Auth) Auth(c *gin.Context) {
	token := c.Request.Header.Get("X-ACCOUNT-TOKEN")
	if token == "" {
		c.JSON(401, userApp.responder.Fail("Unauthorized"))
		c.Abort()
		return
	}
	data, err := userApp.token.Validate(token)
	if err != nil {
		fmt.Println(err)
		c.JSON(401, userApp.responder.Fail(err))
		c.Abort()
		return
	}
	err, dbToken := userApp.userApp.GetTokenByUId(data.Data.TokenId)
	if err != nil {
		c.JSON(401, userApp.responder.Fail("Token was not found"))
		c.Abort()
		return
	}
	if dbToken.AccessTokenExpire.Unix() < time.Now().Unix() {
		c.JSON(401, userApp.responder.Fail("Token is expired"))
		c.Abort()
		return
	}
}
