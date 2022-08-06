package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"user-app/application"
	"user-app/controller"
	"user-app/controller/middleware"
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
	middlewares := middleware.NewMiddleware(app)
	responder := infrastructure.NewResponder()
	authHandlers := controller.NewAuth(app, responder)
	userHandlers := controller.NewUsers(app, responder)

	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())
	authRoutes(router, authHandlers)
	userRoutes(router, userHandlers, middlewares)
	appPort := os.Getenv("API_PORT")
	if appPort == "" {
		appPort = "8888"
	}
	log.Fatal(router.Run(":" + appPort))
}

func authRoutes(router *gin.Engine, handler *controller.Authenticate) {
	router.POST("/auth/login", handler.Login)
}

func userRoutes(router *gin.Engine, users *controller.Users, middleware *middleware.Middleware) {
	router.POST("/users", middleware.AuthMiddleware(), users.SaveUser)
	router.GET("/users", middleware.AuthMiddleware(), users.GetList)
	router.GET("/users/:user_id", users.GetUser)
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
