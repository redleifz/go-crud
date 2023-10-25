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

	//run database
	configs.ConnectDB()

	apiGroup := router.Group("/api/")

	// Mount your routes inside the /api group
	routes.UserRoute(apiGroup) // Example, you can add more routes here
	routes.CelRoutes(apiGroup)

	// router.Run("localhost:6000")'

	router.Run(":" + port)
}
