package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

// InitDB データベース接続を初期化する
func InitDB() {
	var err error

	// 環境変数からデータベースURLを取得
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// データベースに接続
	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	log.Println("Successfully connected to database")
}
