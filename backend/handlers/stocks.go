package handlers

import (
	"fmt"
	"log"
	"math"
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

// 取得用戶所有股票
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

// 新增或更新股票
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

// 刪除股票
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

// ============================
// 📱 WebSocket 即時股價推播
// ============================

var (
	stockClients      = make(map[*websocket.Conn]int) // conn => userID
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
			break // 關閉連線
		}
	}

	stockClientsMutex.Lock()
	delete(stockClients, conn)
	stockClientsMutex.Unlock()
}

// 🧠 每 10 秒抓即時股價並推播
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

func round(val float64) float64 {
	return math.Round(val)
}

// 傳送使用者持股資訊（包含即時價格與損益）
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

		// 台銀計算邏輯
		buyAmount := float64(shares) * avgPrice
		buyFee := round(buyAmount * 0.001425 * 0.35)

		sellAmount := float64(shares) * price
		sellFee := round(sellAmount * 0.001425 * 0.35)
		tax := round(sellAmount * 0.003)

		cost := buyAmount + buyFee
		netSell := sellAmount - sellFee - tax
		profit := round(netSell - cost)

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
		profit := float64(shares) * (price - avgPrice)

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

// 賣出股票處理
func SellStock(c *gin.Context) {
	userID := c.GetInt("user_id")
	var input struct {
		Symbol    string  `json:"symbol"`
		Shares    int     `json:"shares"`
		SellPrice float64 `json:"sell_price"`
		Note      string  `json:"note"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "格式錯誤"})
		return
	}

	// 查詢現有持股
	var currentShares int
	var avgPrice float64
	err := database.DB.QueryRow("SELECT shares, avg_price FROM stocks WHERE user_id = $1 AND symbol = $2",
		userID, input.Symbol).Scan(&currentShares, &avgPrice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "查無持股"})
		return
	}
	if input.Shares <= 0 || input.Shares > currentShares {
		c.JSON(http.StatusBadRequest, gin.H{"error": "賣出股數錯誤"})
		return
	}

	// 損益計算
	buyAmount := float64(input.Shares) * avgPrice
	buyFee := round(buyAmount * 0.001425 * 0.35)

	sellAmount := float64(input.Shares) * input.SellPrice
	sellFee := round(sellAmount * 0.001425 * 0.35)
	tax := round(sellAmount * 0.003)

	cost := buyAmount + buyFee
	netSell := sellAmount - sellFee - tax
	profit := round(netSell - cost)

	// 新增交易紀錄
	_, err = database.DB.Exec(`
		INSERT INTO stock_transactions (user_id, symbol, shares, sell_price, avg_price, realized_profit, note)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		userID, input.Symbol, input.Shares, input.SellPrice, avgPrice, profit, input.Note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "交易紀錄儲存失敗"})
		return
	}

	// 更新持股表
	newShares := currentShares - input.Shares
	if newShares == 0 {
		_, _ = database.DB.Exec("DELETE FROM stocks WHERE user_id = $1 AND symbol = $2", userID, input.Symbol)
	} else {
		_, _ = database.DB.Exec("UPDATE stocks SET shares = $1 WHERE user_id = $2 AND symbol = $3",
			newShares, userID, input.Symbol)
	}

	c.JSON(http.StatusOK, gin.H{"message": "賣出成功", "realized_profit": profit})
}

// 取得交易紀錄
func GetTransactions(c *gin.Context) {
	userID := c.GetInt("user_id")
	rows, err := database.DB.Query(`
		SELECT id, symbol, shares, avg_price, sell_price, realized_profit, note, created_at
		FROM stock_transactions WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "讀取交易紀錄失敗"})
		return
	}
	defer rows.Close()

	type Tx struct {
		ID        int       `json:"id"`
		Symbol    string    `json:"symbol"`
		Shares    int       `json:"shares"`
		AvgPrice  float64   `json:"avg_price"`
		SellPrice float64   `json:"sell_price"`
		Profit    float64   `json:"profit"`
		Note      string    `json:"note"`
		Time      time.Time `json:"created_at"`
	}

	var txs []Tx
	for rows.Next() {
		var t Tx
		if err := rows.Scan(&t.ID, &t.Symbol, &t.Shares, &t.AvgPrice, &t.SellPrice, &t.Profit, &t.Note, &t.Time); err == nil {
			txs = append(txs, t)
		}
	}
	c.JSON(http.StatusOK, gin.H{"transactions": txs})
}

// 總損益摘要 API
func GetStockSummary(c *gin.Context) {
	userID := c.GetInt("user_id")

	var unrealized float64 = 0

	rows, err := database.DB.Query("SELECT symbol, shares, avg_price FROM stocks WHERE user_id = $1", userID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var symbol string
			var shares int
			var avgPrice float64
			if err := rows.Scan(&symbol, &shares, &avgPrice); err == nil {
				price, err := utils.FetchTWSEPrice(symbol)
				if err != nil {
					continue
				}

				// 台銀手續費與交易稅
				buyAmount := float64(shares) * avgPrice
				buyFee := round(buyAmount * 0.001425 * 0.35)
				sellAmount := float64(shares) * price
				sellFee := round(sellAmount * 0.001425 * 0.35)
				tax := round(sellAmount * 0.003)
				cost := buyAmount + buyFee
				netSell := sellAmount - sellFee - tax
				unrealized += round(netSell - cost)
			}
		}
	}

	var realized float64 = 0
	err = database.DB.QueryRow("SELECT COALESCE(SUM(realized_profit),0) FROM stock_transactions WHERE user_id = $1", userID).
		Scan(&realized)
	if err != nil {
		realized = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"unrealized_profit": unrealized,
		"realized_profit":   realized,
		"total_profit":      round(unrealized + realized),
	})
}

// 💰 領取股息 API
func ReceiveDividend(c *gin.Context) {
	userID := c.GetInt("user_id")
	var input struct {
		Symbol string  `json:"symbol"`
		Amount float64 `json:"amount"`
		Note   string  `json:"note"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Amount <= 0 || input.Symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "格式錯誤"})
		return
	}

	note := input.Note + "（股息）"
	_, err := database.DB.Exec(`
		INSERT INTO stock_transactions (user_id, symbol, shares, sell_price, avg_price, realized_profit, note)
		VALUES ($1, $2, 0, 0, 0, $3, $4)`,
		userID, input.Symbol, input.Amount, note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "紀錄股息失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "股息記錄成功"})
}

func ExportTransactionExcel(c *gin.Context) {
	userID := c.GetInt("user_id")

	rows, err := database.DB.Query(`SELECT symbol, shares, avg_price, sell_price, realized_profit, note, created_at
		FROM stock_transactions WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "資料讀取失敗"})
		return
	}
	defer rows.Close()

	f := excelize.NewFile()
	sheet := "交易紀錄"
	f.SetSheetName("Sheet1", sheet)
	headers := []string{"代碼", "股數", "均價", "賣價", "損益", "備註", "時間"}
	for i, h := range headers {
		f.SetCellValue(sheet, fmt.Sprintf("%s1", string(rune('A'+i))), h)
	}

	rowIdx := 2
	for rows.Next() {
		var symbol, note, createdAt string
		var shares int
		var avg, sell, profit float64
		_ = rows.Scan(&symbol, &shares, &avg, &sell, &profit, &note, &createdAt)

		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIdx), symbol)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIdx), shares)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIdx), avg)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIdx), sell)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIdx), profit)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", rowIdx), note)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", rowIdx), createdAt)
		rowIdx++
	}

	filename := fmt.Sprintf("transactions_%d.xlsx", time.Now().Unix())
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	_ = f.Write(c.Writer)
}

func DeleteTransaction(c *gin.Context) {
	userID := c.GetInt("user_id")
	id := c.Param("id")

	_, err := database.DB.Exec("DELETE FROM stock_transactions WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "刪除失敗"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "紀錄已刪除"})
}
