package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func InitPostgres() {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        password TEXT NOT NULL,
        role VARCHAR(20) DEFAULT 'user',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create users table: ", err)
	}

	log.Println("Database initialized successfully")
}

// JWT Middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Register Handler
func Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	_, err := DB.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", input.Name, input.Email, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully!"})
}

// Login Handler
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var storedPassword string
	err := DB.QueryRow("SELECT password FROM users WHERE email = $1", input.Email).Scan(&storedPassword)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": input.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// User Profile Handler
func GetProfile(c *gin.Context) {
	email, _ := c.Get("email")
	var user struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
	}
	err := DB.QueryRow("SELECT id, name, email, created_at FROM users WHERE email = $1", email).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func main() {
	LoadEnv()
	InitPostgres()

	r := gin.Default()

	r.POST("/api/register", Register)
	r.POST("/api/login", Login)

	auth := r.Group("/api").Use(AuthMiddleware())
	auth.GET("/profile", GetProfile)

	r.Static("/assets", "./frontend/assets")
	r.LoadHTMLGlob("./frontend/*.html")
	r.GET("/index", func(c *gin.Context) { c.HTML(200, "index.html", nil) })
	r.GET("/about", func(c *gin.Context) { c.HTML(200, "about.html", nil) })
	r.GET("/projects", func(c *gin.Context) { c.HTML(200, "projects.html", nil) })
	r.GET("/skills", func(c *gin.Context) { c.HTML(200, "skills.html", nil) })
	r.GET("/contact", func(c *gin.Context) { c.HTML(200, "contact.html", nil) })
	r.GET("/auth", func(c *gin.Context) { c.HTML(200, "auth.html", nil) })
	r.GET("/dashboard", func(c *gin.Context) { c.HTML(200, "dashboard.html", nil) })

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
