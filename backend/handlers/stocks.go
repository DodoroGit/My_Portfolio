package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/DodoroGit/My_Portfolio/backend/database"
	"github.com/DodoroGit/My_Portfolio/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xuri/excelize/v2"
)

type Stock struct {
	ID        int       `json:"id"`
	Symbol    string    `json:"symbol"`
	Shares    int       `json:"shares"`
	AvgPrice  float64   `json:"avg_price"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func GetStocks(c *gin.Context) {
	userID := c.GetInt("user_id")
	rows, err := database.DB.Query(`SELECT id, symbol, shares, avg_price, created_at FROM stocks WHERE user_id = $1`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "讀取失敗"})
		return
	}
	defer rows.Close()

	var stocks []Stock
	for rows.Next() {
		var s Stock
		s.UserID = userID
		if err := rows.Scan(&s.ID, &s.Symbol, &s.Shares, &s.AvgPrice, &s.CreatedAt); err == nil {
			stocks = append(stocks, s)
		}
	}
	c.JSON(http.StatusOK, gin.H{"stocks": stocks})
}

func CreateStock(c *gin.Context) {
	userID := c.GetInt("user_id")
	var input struct {
		Symbol   string  `json:"symbol"`
		Shares   int     `json:"shares"`
		AvgPrice float64 `json:"avg_price"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "格式錯誤"})
		return
	}

	_, err := database.DB.Exec(`
		INSERT INTO stocks (user_id, symbol, shares, avg_price)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, symbol)
		DO UPDATE SET shares = EXCLUDED.shares, avg_price = EXCLUDED.avg_price
	`, userID, input.Symbol, input.Shares, input.AvgPrice)
	if err != nil {
		log.Println("❌ SQL 錯誤:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "新增/更新失敗", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "儲存成功"})
}

func DeleteStock(c *gin.Context) {
	userID := c.GetInt("user_id")
	id := c.Param("id")

	_, err := database.DB.Exec("DELETE FROM stocks WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "刪除失敗"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已刪除"})
}

var (
	stockClients      = make(map[*websocket.Conn]int)
	stockClientsMutex sync.Mutex
	stockUpgrader     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func StockSocketHandler(c *gin.Context) {
	userID := c.GetInt("user_id")

	conn, err := stockUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket 連線失敗:", err)
		return
	}
	defer conn.Close()

	stockClientsMutex.Lock()
	stockClients[conn] = userID
	stockClientsMutex.Unlock()

	log.Printf("✅ 使用者 %d 已連線股票 WebSocket\n", userID)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	stockClientsMutex.Lock()
	delete(stockClients, conn)
	stockClientsMutex.Unlock()
}

func StartStockPriceBroadcast() {
	go func() {
		for {
			time.Sleep(10 * time.Second)

			stockClientsMutex.Lock()
			clientsCopy := make(map[*websocket.Conn]int)
			for conn, uid := range stockClients {
				clientsCopy[conn] = uid
			}
			stockClientsMutex.Unlock()

			for conn, uid := range clientsCopy {
				go pushUserStocks(conn, uid)
			}
		}
	}()
}

func pushUserStocks(conn *websocket.Conn, userID int) {
	rows, err := database.DB.Query("SELECT symbol, shares, avg_price FROM stocks WHERE user_id = $1", userID)
	if err != nil {
		return
	}
	defer rows.Close()

	type StockPayload struct {
		Symbol   string  `json:"symbol"`
		Price    float64 `json:"price"`
		Shares   int     `json:"shares"`
		AvgPrice float64 `json:"avg_price"`
		Profit   float64 `json:"profit"`
	}

	for rows.Next() {
		var symbol string
		var shares int
		var avgPrice float64
		if err := rows.Scan(&symbol, &shares, &avgPrice); err != nil {
			continue
		}

		price, err := utils.FetchTWSEPrice(symbol)
		if err != nil {
			continue
		}

		buyCost := float64(shares) * avgPrice * 1.001425
		sellRevenue := float64(shares) * price * (1 - 0.001425 - 0.003)
		profit := sellRevenue - buyCost

		payload := StockPayload{
			Symbol:   symbol,
			Price:    price,
			Shares:   shares,
			AvgPrice: avgPrice,
			Profit:   profit,
		}

		if err := conn.WriteJSON(payload); err != nil {
			log.Println("推播失敗，關閉連線")
			conn.Close()
			stockClientsMutex.Lock()
			delete(stockClients, conn)
			stockClientsMutex.Unlock()
			break
		}
	}
}

func ExportStockExcel(c *gin.Context) {
	userID := c.GetInt("user_id")

	rows, err := database.DB.Query("SELECT symbol, shares, avg_price FROM stocks WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "資料讀取失敗"})
		return
	}
	defer rows.Close()

	f := excelize.NewFile()
	sheet := "持股資料"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"股票代碼", "持股數量", "購入均價", "即時價格", "損益"}
	for i, h := range headers {
		col := string(rune('A'+i)) + "1"
		f.SetCellValue(sheet, col, h)
	}

	rowIndex := 2
	for rows.Next() {
		var symbol string
		var shares int
		var avgPrice float64
		if err := rows.Scan(&symbol, &shares, &avgPrice); err != nil {
			continue
		}

		price, err := utils.FetchTWSEPrice(symbol)
		if err != nil {
			price = 0.0
		}

		buyCost := float64(shares) * avgPrice * 1.001425
		sellRevenue := float64(shares) * price * (1 - 0.001425 - 0.003)
		profit := sellRevenue - buyCost

		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIndex), symbol)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), shares)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIndex), avgPrice)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIndex), price)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIndex), profit)
		rowIndex++
	}

	filename := fmt.Sprintf("stocks_%d.xlsx", time.Now().Unix())
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("File-Name", filename)
	c.Header("Content-Transfer-Encoding", "binary")
	_ = f.Write(c.Writer)
}

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

func GetPortfolioSummary(c *gin.Context) {
	userID := c.GetInt("user_id")

	rows, err := database.DB.Query("SELECT symbol, shares, avg_price FROM stocks WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "資料讀取失敗"})
		return
	}
	defer rows.Close()

	var totalCost float64 = 0
	var totalValue float64 = 0

	for rows.Next() {
		var symbol string
		var shares int
		var avgPrice float64
		if err := rows.Scan(&symbol, &shares, &avgPrice); err != nil {
			continue
		}

		price, err := utils.FetchTWSEPrice(symbol)
		if err != nil {
			continue
		}

		cost := float64(shares) * avgPrice * 1.001425
		value := float64(shares) * price * (1 - 0.001425 - 0.003)

		totalCost += cost
		totalValue += value
	}

	c.JSON(http.StatusOK, gin.H{
		"total_cost":     fmt.Sprintf("%.2f", totalCost),
		"total_value":    fmt.Sprintf("%.2f", totalValue),
		"unrealized_pnl": fmt.Sprintf("%.2f", totalValue-totalCost),
	})
}
