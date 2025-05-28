// handlers/stock_history.go
package handlers

import (
	"net/http"

	"github.com/DodoroGit/My_Portfolio/backend/utils"
	"github.com/gin-gonic/gin"
)

func GetStockHistory(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請提供股票代碼"})
		return
	}

	history, err := utils.FetchTWSEHistory(symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得歷史資料"})
		return
	}

	c.JSON(http.StatusOK, history)
}
