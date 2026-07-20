package main

import (
	"net/http"
	"os"

	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	// Auto migrate
	err := config.DB.AutoMigrate(&models.Task{})
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API Running",
		})
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}