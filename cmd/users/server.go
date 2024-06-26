package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"user-app/internal/application"
	http2 "user-app/internal/http"
	middleware2 "user-app/internal/http/middleware"
	"user-app/internal/repository"
	"user-app/pkg/utils"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("no env gotten")
	}
}

func main() {
	repositories := getRepositories()
	repositories.Automigrate() //todo this is for dev environment
	app := application.NewUserApp(repositories.User, repositories.Token)
	token := utils.NewToken(utils.GetFileBytes(os.Getenv("PRIVATE_KEY_FILE_PATH")), utils.GetFileBytes(os.Getenv("PUBLIC_KEY_FILE_PATH")))
	responder := utils.NewResponder()
	authMiddleware := middleware2.NewAuth(app, responder, token)
	authHandlers := http2.NewLoginHandler(app, responder, token)
	userHandlers := http2.NewUserHandler(app, responder)

	router := gin.Default()
	router.Use(middleware2.Cors())
	router.Static("/static", "./static")
	authRoutes(router, authHandlers)
	userRoutes(router, userHandlers, authMiddleware)
	appPort := os.Getenv("API_PORT")
	if appPort == "" {
		log.Fatal("API_PORT variable is not set in the .env file")
	}

	if os.Getenv("APP_ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, responder.Fail("404 not found"))
	})

	log.Fatal(router.Run(":" + appPort))
}

func authRoutes(router *gin.Engine, handler *http2.LoginHandler) {
	auth := router.Group("/v1/auth")
	auth.POST("/login", handler.Login)
}

func userRoutes(router *gin.Engine, users *http2.UserHandler, middleware *middleware2.Auth) {
	userGroup := router.Group("/v1/users")
	userGroup.POST("/", middleware.Auth, users.SaveUser)
	userGroup.GET("/", middleware.Auth, users.GetList)
	userGroup.GET("/:user_id", middleware.Auth, users.GetUser)
}

func getRepositories() *repository.Repositories {
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	repositories, err := repository.NewRepositories(user, password, port, host, dbname)
	if err != nil {
		panic(err)
	}

	return repositories
}
