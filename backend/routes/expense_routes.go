package routes

import (
	"github.com/DodoroGit/My_Portfolio/backend/handlers"
	"github.com/DodoroGit/My_Portfolio/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func ExpenseRoutes(r *gin.Engine) {
	expense := r.Group("/api/expense")
	expense.Use(middlewares.AuthMiddleware())
	{
		expense.POST("/", handlers.CreateExpense)
		expense.GET("/", handlers.GetExpenses)
		expense.POST("/upload", handlers.UploadExcel) // 後續擴充
		expense.GET("/export", handlers.ExportExcel)  // 後續擴充
	}
}
