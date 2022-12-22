package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"user-app/application"
	"user-app/http"
	"user-app/http/middleware"
	"user-app/infrastructure"
	"user-app/infrastructure/repository"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("no env gotten")
	}
}

func main() {
	repositories := getRepositories()
	defer repositories.Close()
	repositories.Automigrate() //todo this is for dev environment
	app := application.NewUserApp(repositories.User, repositories.Token)
	responder := infrastructure.NewResponder()
	authMiddleware := middleware.NewAuth(app, responder)
	authHandlers := http.NewLoginHandler(app, responder)
	userHandlers := http.NewUserHandler(app, responder)

	router := gin.Default()
	router.Use(middleware.Cors())
	authRoutes(router, authHandlers)
	userRoutes(router, userHandlers, authMiddleware)
	appPort := os.Getenv("API_PORT")
	if appPort == "" {
		log.Fatal("API_PORT variable is not set in the .env file")
	}
	log.Fatal(router.Run(":" + appPort))
}

func authRoutes(router *gin.Engine, handler *http.LoginHandler) {
	auth := router.Group("/v1/auth")
	auth.POST("/login", handler.Login)
}

func userRoutes(router *gin.Engine, users *http.UserHandler, middleware *middleware.Auth) {
	userGroup := router.Group("/v1/users")
	userGroup.POST("/", middleware.Auth, users.SaveUser)
	userGroup.GET("/", middleware.Auth, users.GetList)
	userGroup.GET("/:user_id", middleware.Auth, users.GetUser)
}

func getRepositories() *repository.Repositories {
	dbDriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	repositories, err := repository.NewRepositories(dbDriver, user, password, port, host, dbname)
	if err != nil {
		panic(err)
	}

	return repositories
}
