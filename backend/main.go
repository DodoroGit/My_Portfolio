package main

import (
	"log"

	"github.com/DodoroGit/My_Portfolio/backend/database"
	"github.com/DodoroGit/My_Portfolio/backend/handlers"
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

	//routes.WebRoutes(r) 透過Frontend Nginx反向代理替換
	routes.RegisterRoutes(r)
	// WebSocket 路由
	handlers.StartStockPriceBroadcast()

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
