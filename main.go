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
		log.Println("no env gotten")
	}
}

func main() {
	var app application.UserApp
	repositories := getRepositories()
	defer repositories.Close()
	repositories.Automigrate()
	app = application.UserApp{UserRepo: *repositories.User, AccessTokenRepo: *repositories.Token}
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
		appPort = "8081"
	}
	log.Fatal(router.Run(":" + appPort))
}

func authRoutes(router *gin.Engine, handler *http.LoginHandler) {
	router.POST("/auth/login", handler.Login)
}

func userRoutes(router *gin.Engine, users *http.UserHandler, middleware *middleware.Auth) {
	router.POST("/v1/users", middleware.Auth(), users.SaveUser)
	router.GET("/v1/users", middleware.Auth(), users.GetList)
	router.GET("/v1/users/:user_id", users.GetUser)
}

func getRepositories() *repository.Repositories {
	dbDriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	repositories, err := repository.CreateRepositories(dbDriver, user, password, port, host, dbname)
	if err != nil {
		panic(err)
	}

	return repositories
}
