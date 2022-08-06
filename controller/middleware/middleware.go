package middleware

import (
	"github.com/gin-gonic/gin"
	"user-app/application"
)

type Middleware struct {
	userApp application.UserApp
}

func NewMiddleware(userApp application.UserApp) *Middleware {
	return &Middleware{userApp}
}

func (userApp *Middleware) AuthMiddleware() gin.HandlerFunc {
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
		c.Next()
	}
}

func (userApp *Middleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
