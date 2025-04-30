package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB() {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "startios" // по умолчанию — основная база
	}

	dsn := fmt.Sprintf("host=localhost user=postgres password=7292 dbname=%s port=5432 sslmode=disable", dbName)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Ошибка подключения к БД:", err)
	}

	log.Println("✅ Database connected:", dbName)
}
