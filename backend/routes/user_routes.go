package routes

import (
	"github.com/DodoroGit/My_Portfolio/backend/handlers"
	"github.com/DodoroGit/My_Portfolio/backend/middlewares"

	"github.com/gin-gonic/gin"
)

// UserRoutes 設定使用者 API 路由
func UserRoutes(r *gin.Engine) {
	user := r.Group("/api/user")
	user.Use(middlewares.AuthMiddleware()) // 保護這些路由
	{
		user.GET("/profile", handlers.GetProfile)
		user.PUT("/profile", handlers.UpdateProfile)
	}

	admin := r.Group("/api/admin")
	admin.Use(middlewares.AuthMiddleware()) // 保護管理員 API
	{
		admin.GET("/users", handlers.GetAllUsers)
	}
}
