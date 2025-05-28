package routes

import (
	"github.com/DodoroGit/My_Portfolio/backend/handlers"
	"github.com/DodoroGit/My_Portfolio/backend/middlewares"
	"github.com/gin-gonic/gin"
)

// StockRoutes 設定股票追蹤功能的 API 路由
func StockRoutes(r *gin.Engine) {
	stock := r.Group("/api/stocks")
	stock.Use(middlewares.AuthMiddleware())
	{
		stock.GET("/", handlers.GetStocks)                      // 取得所有持股資料
		stock.POST("/", handlers.CreateStock)                   // 新增或更新持股資料
		stock.DELETE("/:id", handlers.DeleteStock)              // 刪除持股資料
		stock.GET("/export", handlers.ExportStockExcel)         // 匯出持股資料為 Excel 檔案
		stock.GET("/history/:symbol", handlers.GetStockHistory) // 取得股票歷史資料
	}

	// ⭐️ 加入 WebSocket 路由
	ws := r.Group("/ws/stocks")
	ws.Use(middlewares.AuthMiddleware())
	{
		ws.GET("/", handlers.StockSocketHandler)
	}
}
