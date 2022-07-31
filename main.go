package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"user-app/controller"
	"user-app/controller/middleware"
	"user-app/infrastructure/persistence"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

func main() {
	repositories := getRepositories()
	defer repositories.Close()
	repositories.Automigrate()
	authHandlers := controller.NewAuthenticate(repositories.User)
	userHandlers := controller.NewUsers(repositories.User)

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	authRoutes(router, authHandlers)
	userRoutes(router, userHandlers)
	appPort := os.Getenv("API_PORT")
	if appPort == "" {
		appPort = "8888"
	}
	log.Fatal(router.Run(":" + appPort))
}

func authRoutes(router *gin.Engine, handler *controller.Authenticate) {
	router.POST("/auth/login", handler.Login)
	router.POST("/auth/logout", handler.Logout)
}

func userRoutes(router *gin.Engine, users *controller.Users) {
	router.POST("/users", middleware.AuthMiddleware(), users.SaveUser)
	router.GET("/users", users.GetList)
	router.GET("/users/:user_id", users.GetUser)
}

func getRepositories() *persistence.Repositories {
	dbDriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	repositories, err := persistence.CreateRepositories(dbDriver, user, password, port, host, dbname)
	if err != nil {
		panic(err)
	}

	return repositories
}
