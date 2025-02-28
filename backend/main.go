package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // 🚀 這行很重要，確保 PostgreSQL 驅動被加載
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
	projects := []map[string]string{
		{"name": "個人作品集", "description": "展示我的個人專案", "url": "projects.html"},
		{"name": "E-commerce 平台", "description": "一個簡單的線上購物網站", "url": "#"},
		{"name": "部落格系統", "description": "基於 Gin 框架開發的部落格平台", "url": "#"},
	}

	c.JSON(200, gin.H{"projects": projects})
}

func main() {
	LoadEnv()
	InitPostgres()

	r := gin.Default()

	r.Static("/static", "./frontend")

	r.GET("/users", GetUsers)
	r.GET("/projects", GetProjects)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatal("Unable to start server: ", err)
	}

}
