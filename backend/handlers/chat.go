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

// 管理所有連線中的使用者
var clients = make(map[*websocket.Conn]string) // 連線 -> 使用者名稱
var clientsMutex = sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	UserID    int       `json:"user_id"` // ⭐️新增 user_id
	UserName  string    `json:"user_name"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// WebSocket 聊天處理
func ChatHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	userName := c.GetString("user_name")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket 升級失敗:", err)
		return
	}
	defer conn.Close()

	clientsMutex.Lock()
	clients[conn] = userName
	clientsMutex.Unlock()

	// 初始發送最近 500 則訊息
	rows, err := database.DB.Query("SELECT user_name, content, created_at FROM messages ORDER BY created_at DESC LIMIT 500")
	if err == nil {
		defer rows.Close()
		var history []Message
		for rows.Next() {
			var msg Message
			if err := rows.Scan(&msg.UserName, &msg.Content, &msg.Timestamp); err == nil {
				history = append([]Message{msg}, history...) // 倒序
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
			incoming.Content = incoming.Content[:500] // 限制最大長度
		}

		incoming.UserID = userID
		incoming.UserName = userName
		incoming.Timestamp = time.Now()

		// 保存到資料庫
		_, err = database.DB.Exec("INSERT INTO messages (user_id, user_name, content) VALUES ($1, $2, $3)", userID, userName, incoming.Content)
		if err != nil {
			log.Println("儲存訊息到資料庫失敗:", err)
		}

		// 廣播給所有連線
		broadcastMessage(incoming)
	}

	// 使用者離開，清理連線
	clientsMutex.Lock()
	delete(clients, conn)
	clientsMutex.Unlock()
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
