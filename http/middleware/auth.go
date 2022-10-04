package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"time"
	"user-app/application"
	"user-app/infrastructure"
	jwt_token "user-app/infrastructure/jwt-token"
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
		privateKey, err := os.ReadFile("private_key.pem")
		if err != nil {
			log.Fatalln(err)
		}
		publicKey, err := os.ReadFile("public_key.pem")
		if err != nil {
			log.Fatalln(err)
		}
		jwt := jwt_token.NewJWT(privateKey, publicKey)
		_, err = jwt.Validate(token)
		if err != nil {
			fmt.Println(err)
			c.JSON(401, userApp.responder.Fail(err))
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
			c.JSON(401, userApp.responder.Fail("jwt-token is expired"))
			c.Abort()
			return
		}

		c.Next()
	}
}
