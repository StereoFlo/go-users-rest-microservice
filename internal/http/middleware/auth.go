package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"user-app/internal/application"
	infrastructure2 "user-app/internal/infrastructure"
)

type Auth struct {
	userApp   *application.UserApp
	responder *infrastructure2.Responder
}

func NewAuth(userApp *application.UserApp, responder *infrastructure2.Responder) *Auth {
	return &Auth{userApp, responder}
}

func (userApp *Auth) Auth(c *gin.Context) {
	token := c.Request.Header.Get("X-ACCOUNT-TOKEN")
	if token == "" {
		c.JSON(401, userApp.responder.Fail("Unauthorized"))
		c.Abort()
		return
	}
	jwt := infrastructure2.NewToken()
	data, err := jwt.Validate(token)
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
