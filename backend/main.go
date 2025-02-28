package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func InitPostgres() {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
}

func GetUsers(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get all users"})
}

func GetProjects(c *gin.Context) {
	c.File("./frontend/index.html")
}

func main() {
	LoadEnv()
	InitPostgres()

	r := gin.Default()

	// 🚀 新增靜態文件伺服，讓 Gin 服務 frontend 資料夾的靜態文件
	r.Static("/static", "./frontend")

	// 🚀 設定首頁 ("/") 轉向 index.html
	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	// API 路由
	r.GET("/users", GetUsers)
	r.GET("/projects", GetProjects)

	// 讓 Gin 監聽 0.0.0.0:8080
	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
