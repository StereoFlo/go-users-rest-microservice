package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
	"user-app/application"
	"user-app/infrastructure"
)

type Auth struct {
	userApp   application.UserApp
	responder *infrastructure.Responder
}

func NewAuth(userApp application.UserApp, responder *infrastructure.Responder) *Auth {
	return &Auth{userApp, responder}
}

func (userApp *Auth) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("X-ACCOUNT-TOKEN")
		if token == "" {
			c.JSON(401, userApp.responder.Fail("unauthorized"))
			c.Abort()
			return
		}

		if !isValidUUID(token) {
			c.JSON(401, userApp.responder.Fail("token is wrong"))
			c.Abort()
			return
		}

		bdToken, err := userApp.userApp.UserRepo.GetUserByAccessToken(token)
		if err != nil {
			c.JSON(401, userApp.responder.Fail(err))
			c.Abort()
			return
		}

		if bdToken.AccessTokenExpire.Unix() < time.Now().Unix() {
			c.JSON(401, userApp.responder.Fail("token is expired"))
			c.Abort()
			return
		}

		c.Next()
	}
}
func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
