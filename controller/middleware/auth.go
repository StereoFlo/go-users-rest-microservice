package middleware

import (
	"github.com/gin-gonic/gin"
	"time"
	"user-app/application"
)

type Auth struct {
	userApp application.UserApp
}

func NewAuth(userApp application.UserApp) *Auth {
	return &Auth{userApp}
}

func (userApp *Auth) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("X-ACCOUNT-TOKEN")
		if token == "" {
			c.JSON(401, gin.H{
				"meta": gin.H{
					"success": false,
				},
				"data": "unauthorized",
			})
			c.Abort()
		}
		bdToken, _ := userApp.userApp.UserRepo.GetUserByAccessToken(token)
		if bdToken == nil {
			c.JSON(401, gin.H{
				"meta": gin.H{
					"success": false,
				},
				"data": "token is wrong",
			})
			c.Abort()
		}

		if bdToken.AccessTokenExpire.Unix() > time.Now().Unix() {
			c.JSON(401, gin.H{
				"meta": gin.H{
					"success": false,
				},
				"data": "token is expired",
			})
			c.Abort()
		}

		c.Next()
	}
}
