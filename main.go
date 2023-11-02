package main

import (
	"go-crud/configs"
	"go-crud/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env fileddd ")
	}

	// Get the port number from the environment variables
	port := os.Getenv("APP_PORT")
	router := gin.Default()

	// Use CORS middleware to enable Cross-Origin Resource Sharing
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Max-Age", "600")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	//run database
	configs.ConnectDB()

	apiGroup := router.Group("/api/")

	// Mount your routes inside the /api group
	routes.UserRoute(apiGroup) // Example, you can add more routes here
	routes.CelRoutes(apiGroup)

	// router.Run("localhost:6000")'

	router.Run(":" + port)
}

// CORSMiddleware creates a CORS middleware handler
