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
	c.HTML(200, "index.html", nil)
}

func main() {
	LoadEnv()
	InitPostgres()

	r := gin.Default()

	r.Static("/assets", "/home/ec2-user/My_Portfolio/frontend/assets")
	r.LoadHTMLGlob("/home/ec2-user/My_Portfolio/frontend/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/about", func(c *gin.Context) {
		c.HTML(200, "about.html", nil)
	})

	r.GET("/projects", func(c *gin.Context) {
		c.HTML(200, "projects.html", nil)
	})

	r.GET("/skills", func(c *gin.Context) {
		c.HTML(200, "skills.html", nil)
	})

	r.GET("/contact", func(c *gin.Context) {
		c.HTML(200, "contact.html", nil)
	})

	r.GET("/users", GetUsers)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatal("Unable to start server: ", err)
	}

}
