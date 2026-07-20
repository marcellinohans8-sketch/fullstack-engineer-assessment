package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/tasks", controllers.CreateTask)
api.GET("/tasks", controllers.GetTasks)
api.PUT("/tasks/:id", controllers.UpdateTask)
api.DELETE("/tasks/:id", controllers.DeleteTask)
	}
}