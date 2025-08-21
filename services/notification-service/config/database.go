package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"notification-service/internal/model"
	"os"
	"strconv"
	"time"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Validate required env vars
	missingVars := []string{}
	if dbUser == "" {
		missingVars = append(missingVars, "DB_USER")
	}
	if dbPass == "" {
		missingVars = append(missingVars, "DB_PASS")
	}
	if dbHost == "" {
		missingVars = append(missingVars, "DB_HOST")
	}
	if dbPort == "" {
		missingVars = append(missingVars, "DB_PORT")
	}
	if dbName == "" {
		missingVars = append(missingVars, "DB_NAME")
	}

	if len(missingVars) > 0 {
		log.Fatalf("❌ Missing required environment variables: %v", missingVars)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	// Add retry logic with timeout
	var db *gorm.DB
	maxRetries := getEnvAsInt("DB_MAX_RETRIES", 5)
	retryDelay := time.Duration(getEnvAsInt("DB_RETRY_DELAY_SEC", 3)) * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		log.Printf("Connection attempt %d failed: %v\n", i+1, err)
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to connect after %d attempts: %v", maxRetries, err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get generic database object: " + err.Error())
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db
	err = db.AutoMigrate(model.Notification{})
	if err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}
	log.Println("✅ Connected to MySQL successfully!")
}

func getEnvAsInt(name string, defaultVal int) int {
	if valStr := os.Getenv(name); valStr != "" {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}
