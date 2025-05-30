package routes

import (
	"github.com/DodoroGit/My_Portfolio/backend/handlers"
	"github.com/DodoroGit/My_Portfolio/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// Auth Routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// User Routes
	user := r.Group("/api/user")
	user.Use(middlewares.AuthMiddleware())
	{
		user.GET("/profile", handlers.GetProfile)
		user.PUT("/profile", handlers.UpdateProfile)
	}

	admin := r.Group("/api/admin")
	admin.Use(middlewares.AuthMiddleware())
	{
		admin.GET("/users", handlers.GetAllUsers)
		admin.GET("/pending-users", handlers.GetPendingUsers)
		admin.POST("/approve-user", handlers.ApproveUser)
	}

	// Expense Routes
	expense := r.Group("/api/expense")
	expense.Use(middlewares.AuthMiddleware())
	{
		expense.POST("/", handlers.CreateExpense)
		expense.GET("/", handlers.GetExpenses)
		expense.POST("/upload", handlers.UploadExcel)
		expense.GET("/export", handlers.ExportExcel)
	}

	// Stock Routes
	stock := r.Group("/api/stocks")
	stock.Use(middlewares.AuthMiddleware())
	{
		stock.GET("/", handlers.GetStocks)
		stock.POST("/", handlers.CreateStock)
		stock.DELETE("/:id", handlers.DeleteStock)
		stock.GET("/export", handlers.ExportStockExcel)
		stock.GET("/history/:symbol", handlers.GetStockHistory)
		stock.GET("/summary", handlers.GetPortfolioSummary)
		stock.POST("/sell", handlers.SellStock)
	}

	wsStocks := r.Group("/ws/stocks")
	wsStocks.Use(middlewares.AuthMiddleware())
	{
		wsStocks.GET("/", handlers.StockSocketHandler)
	}

	// Chat Routes
	wsChat := r.Group("/ws")
	wsChat.Use(middlewares.AuthMiddleware())
	{
		wsChat.GET("/chat", handlers.ChatHandler)
	}

	apiChat := r.Group("/api/chat")
	apiChat.Use(middlewares.AuthMiddleware())
	{
		apiChat.POST("/clear", handlers.ClearChatHandler)
	}
}
