package main

import (
	"log"
	"os"

	"github.com/DodoroGit/My_Portfolio/backend/database"
	"github.com/DodoroGit/My_Portfolio/backend/handlers"
	"github.com/DodoroGit/My_Portfolio/backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化資料庫
	database.InitPostgres()

	// 初始化 Gin
	r := gin.Default()

	// 註冊 RESTful 路由
	routes.RegisterRoutes(r)

	// 啟動 WebSocket（股票功能）
	handlers.StartStockPriceBroadcast()

	// 啟動伺服器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
