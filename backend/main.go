package main

import (
	"log"

	"github.com/DodoroGit/My_Portfolio/backend/config"
	"github.com/DodoroGit/My_Portfolio/backend/db"
	"github.com/DodoroGit/My_Portfolio/backend/route"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration
	config.LoadEnv()
	// Initialize database
	db.InitPostgres()
	// Create a Gin router
	r := gin.Default()

	// Setup routes
	route.SetupRoutes(r)
	// Run the server
	if err := r.Run(); err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
