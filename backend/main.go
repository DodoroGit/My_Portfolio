package main

import (
	"log"

	"github.com/DodoroGit/My_Portfolio/backend/database"
	"github.com/DodoroGit/My_Portfolio/backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	//LoadEnv()
	database.InitPostgres()

	r := gin.Default()

	routes.AuthRoutes(r)
	routes.WebRoutes(r)
	routes.UserRoutes(r)
	routes.ChatRoutes(r)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
