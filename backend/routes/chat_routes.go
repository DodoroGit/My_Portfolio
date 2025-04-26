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

	// ⭐ 新增API路由：清除聊天紀錄
	api := r.Group("/api/chat")
	api.Use(middlewares.AuthMiddleware()) // 一樣保護，必須登入
	{
		api.POST("/clear", handlers.ClearChatHandler)
	}
}
