package handlers

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/DodoroGit/My_Portfolio/backend/database"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]string)
var clientsMutex = sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	Type      string    `json:"type"` // ⭐ 新增：訊息類型 (system / message)
	UserID    int       `json:"user_id"`
	UserName  string    `json:"user_name"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func ChatHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	userName := c.GetString("user_name")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket升級失敗:", err)
		return
	}
	defer conn.Close()

	clientsMutex.Lock()
	clients[conn] = userName
	clientsMutex.Unlock()

	// ⭐ 進入時廣播加入系統訊息
	broadcastSystemMessage(userName + " 已加入聊天室")

	// 發送歷史訊息
	rows, err := database.DB.Query("SELECT user_id, user_name, content, created_at FROM messages ORDER BY created_at DESC LIMIT 300")
	if err == nil {
		defer rows.Close()
		var history []Message
		for rows.Next() {
			var msg Message
			if err := rows.Scan(&msg.UserID, &msg.UserName, &msg.Content, &msg.Timestamp); err == nil {
				msg.Type = "message" // ⭐ 歷史訊息都是普通訊息
				history = append([]Message{msg}, history...)
			}
		}
		for _, msg := range history {
			conn.WriteJSON(msg)
		}
	}

	for {
		var incoming Message
		err := conn.ReadJSON(&incoming)
		if err != nil {
			log.Println("讀取訊息失敗:", err)
			break
		}

		if len(incoming.Content) > 500 {
			incoming.Content = incoming.Content[:500]
		}

		incoming.UserID = userID
		incoming.UserName = userName
		incoming.Timestamp = time.Now()
		incoming.Type = "message"

		// 存入資料庫
		_, err = database.DB.Exec("INSERT INTO messages (user_id, user_name, content) VALUES ($1, $2, $3)", userID, userName, incoming.Content)
		if err != nil {
			log.Println("儲存訊息失敗:", err)
		}

		// ⭐ 插入後刪除超過300筆的舊訊息
		_, _ = database.DB.Exec(`
			DELETE FROM messages
			WHERE id NOT IN (
				SELECT id FROM messages ORDER BY created_at DESC LIMIT 300
			)
		`)

		// 廣播
		broadcastMessage(incoming)
	}

	clientsMutex.Lock()
	delete(clients, conn)
	clientsMutex.Unlock()

	// ⭐ 離開時廣播離開系統訊息
	broadcastSystemMessage(userName + " 已離開聊天室")
}

func broadcastMessage(message Message) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for conn := range clients {
		err := conn.WriteJSON(message)
		if err != nil {
			log.Println("廣播訊息失敗:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}

// ⭐ 新增：廣播系統提示訊息
func broadcastSystemMessage(content string) {
	msg := Message{
		Type:      "system",
		Content:   content,
		Timestamp: time.Now(),
	}
	broadcastMessage(msg)
}

// ⭐ 新增：清除所有聊天紀錄API
func ClearChatHandler(c *gin.Context) {
	_, err := database.DB.Exec("DELETE FROM messages")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法清空聊天紀錄"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "聊天紀錄已清空"})
}
