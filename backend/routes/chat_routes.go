package routes

import (
	"github.com/DodoroGit/My_Portfolio/backend/handlers"
	"github.com/DodoroGit/My_Portfolio/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func ChatRoutes(r *gin.Engine) {
	chat := r.Group("/ws")
	chat.Use(middlewares.AuthMiddleware()) // 必須登入
	{
		chat.GET("/chat", handlers.ChatHandler)
	}
}
