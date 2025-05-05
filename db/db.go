package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func InitDB() {
	time.Sleep(5 * time.Second)
	
	// Чтение переменных окружения
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Значения по умолчанию, если переменные не заданы (при запуске без Docker)
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbPassword == "" {
		dbPassword = "7292"
	}
	if dbName == "" {
		dbName = "startios"
	}

	// Формирование строки подключения
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)

	// Подключение к базе
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Ошибка подключения к БД:", err)
	}

	log.Println("✅ Database connected:", dbName)
}

//package db
//
//import (
//	"fmt"
//	"gorm.io/driver/postgres"
//	"gorm.io/gorm"
//	"log"
//	"os"
//)
//
//var DB *gorm.DB
//
//func InitDB() {
//	dbName := os.Getenv("DB_NAME")
//	if dbName == "" {
//		dbName = "startios" // по умолчанию — основная база
//	}
//
//	dsn := fmt.Sprintf("host=localhost user=postgres password=7292 dbname=%s port=5432 sslmode=disable", dbName)
//	var err error
//	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
//	if err != nil {
//		log.Fatal("❌ Ошибка подключения к БД:", err)
//	}
//
//	log.Println("✅ Database connected:", dbName)
//}
