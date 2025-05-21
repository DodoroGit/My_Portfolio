package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/DodoroGit/My_Portfolio/backend/database"
	"github.com/DodoroGit/My_Portfolio/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type CreateExpenseInput struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount" binding:"required"`
	Note     string  `json:"note"`
	SpentAt  string  `json:"spent_at" binding:"required"` // ISO date string
}

func CreateExpense(c *gin.Context) {
	var input CreateExpenseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt("user_id")
	date, err := time.Parse("2006-01-02", input.SpentAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式錯誤"})
		return
	}

	_, err = database.DB.Exec(`
		INSERT INTO expenses (user_id, category, amount, note, spent_at)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, input.Category, input.Amount, input.Note, date)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "儲存失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "新增成功"})
}

func GetExpenses(c *gin.Context) {
	userID := c.GetInt("user_id")

	rows, err := database.DB.Query(`
		SELECT id, category, amount, note, spent_at, created_at
		FROM expenses
		WHERE user_id = $1
		ORDER BY spent_at DESC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢失敗"})
		return
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		e.UserID = userID
		if err := rows.Scan(&e.ID, &e.Category, &e.Amount, &e.Note, &e.SpentAt, &e.CreatedAt); err == nil {
			expenses = append(expenses, e)
		}
	}

	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}

func ExportExcel(c *gin.Context) {
	userID := c.GetInt("user_id")
	rows, err := database.DB.Query(`
		SELECT category, amount, note, spent_at
		FROM expenses
		WHERE user_id = $1
		ORDER BY spent_at DESC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "匯出失敗"})
		return
	}
	defer rows.Close()

	file := excelize.NewFile()
	sheet := "Expenses"
	file.NewSheet(sheet)

	headers := []string{"日期", "類別", "金額", "備註"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		file.SetCellValue(sheet, cell, header)
	}

	rowNum := 2
	for rows.Next() {
		var category, note string
		var amount float64
		var spentAt time.Time
		rows.Scan(&category, &amount, &note, &spentAt)

		file.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), spentAt.Format("2006-01-02"))
		file.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), category)
		file.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), amount)
		file.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), note)
		rowNum++
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=expenses.xlsx")
	c.Header("File-Name", "expenses.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	_ = file.Write(c.Writer)
}

func UploadExcel(c *gin.Context) {
	userID := c.GetInt("user_id")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請提供 Excel 檔案"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法讀取檔案"})
		return
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無法解析 Excel 檔案"})
		return
	}

	rows, err := f.GetRows("Expenses")
	if err != nil || len(rows) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Excel 資料格式錯誤或為空"})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "資料庫交易失敗"})
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM expenses WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "刪除舊資料失敗"})
		return
	}

	stmt, err := tx.Prepare("INSERT INTO expenses (user_id, category, amount, note, spent_at) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "準備寫入失敗"})
		return
	}
	defer stmt.Close()

	for i, row := range rows[1:] {
		if len(row) < 4 {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("第 %d 行資料不完整", i+2)})
			return
		}

		date, err := time.Parse("2006-01-02", row[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("第 %d 行日期格式錯誤", i+2)})
			return
		}

		amount, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("第 %d 行金額格式錯誤", i+2)})
			return
		}

		_, err = stmt.Exec(userID, row[1], amount, row[3], date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "資料寫入失敗"})
			return
		}
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交資料失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "上傳並更新成功"})
}
