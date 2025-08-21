package main

import (
	"api-gateway/internal/i18n"
	"api-gateway/internal/middleware"
	"api-gateway/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

func initEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}
}

func initI18n() {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "..", "locales")
	i18n.SetLocalesDir(basePath)

	if err := i18n.LoadAllLanguages([]string{"en", "vi"}); err != nil {
		log.Fatalf("failed to load languages: %v", err)
	}
}

func initRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.I18nMiddleware())
	r.Use(middleware.TranslateMiddleware())

	// Rate limit
	rateLimit := 100 // default
	if rlStr := os.Getenv("RATE_LIMIT"); rlStr != "" {
		if rl, err := strconv.Atoi(rlStr); err == nil && rl > 0 {
			rateLimit = rl
		}
	}
	r.Use(middleware.RateLimitMiddleware(rateLimit))
	r.Use(middleware.I18nMiddleware())

	// Test route
	r.GET("/hello", func(c *gin.Context) {
		T := c.MustGet("T").(func(string) string)
		c.JSON(200, gin.H{
			"message": T("hello"),
		})
	})

	// Register other routes
	routes.RegisterRoutes(r)

	return r
}

func main() {
	initEnv()
	initI18n()
	r := initRouter()

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
