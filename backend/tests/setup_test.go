package tests

import (
	"testing"

	"backend/config"
	"backend/models"
	"backend/routes"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(&models.Task{}); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	redisServer, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start test redis: %v", err)
	}

	config.DB = db
	config.RedisClient = redis.NewClient(&redis.Options{Addr: redisServer.Addr()})

	t.Cleanup(func() {
		_ = config.RedisClient.Close()
		redisServer.Close()
	})

	router := gin.New()
	routes.SetupRoutes(router)

	return router
}
