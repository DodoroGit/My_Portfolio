package handlers

import (
	"net/http"

	"github.com/DodoroGit/My_Portfolio/backend/database"
	"github.com/DodoroGit/My_Portfolio/backend/models"

	"github.com/gin-gonic/gin"
)

// GetProfile 取得使用者個人資料
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	err := database.DB.QueryRow("SELECT id, name, email, role, created_at FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

type UpdateProfileInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateProfile 更新使用者個人資料
func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3",
		input.Name, input.Email, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// GetAllUsers 取得所有使用者 (需 Admin 權限)
func GetAllUsers(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	rows, err := database.DB.Query("SELECT id, name, email, role, created_at FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning user"})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetPendingUsers 獲取待審核的用戶列表
func GetPendingUsers(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "權限不足"})
		return
	}

	rows, err := database.DB.Query("SELECT id, name, email, role, status, created_at FROM users WHERE status = 'pending'")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "獲取待審核用戶失敗"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Status, &user.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "處理用戶數據時出錯"})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{"pending_users": users})
}

type ApprovalInput struct {
	UserID int    `json:"user_id" binding:"required"`
	Action string `json:"action" binding:"required,oneof=approve reject"` // approve 或 reject
}

// ApproveUser 審核用戶註冊
func ApproveUser(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "權限不足"})
		return
	}

	var input ApprovalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := "approved"
	if input.Action == "reject" {
		status = "rejected"
	}

	_, err := database.DB.Exec("UPDATE users SET status = $1 WHERE id = $2", status, input.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用戶狀態失敗"})
		return
	}

	message := "已批准用戶註冊"
	if input.Action == "reject" {
		message = "已拒絕用戶註冊"
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}
