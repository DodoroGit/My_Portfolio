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
	c.JSON(200, gin.H{"message": "Get all projects"})
}

func main() {
	LoadEnv()
	InitPostgres()

	r := gin.Default()
	r.GET("/users", GetUsers)
	r.GET("/projects", GetProjects)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatal("Unable to start server: ", err)
	}

}
