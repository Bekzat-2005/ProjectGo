package db

import (
	"log"
	"projectGolang/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=7292 dbname=startios port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	log.Println("Database connected!")

	// Миграции
	err = DB.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{})
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	log.Println("Миграция успешно завершена!")
}
