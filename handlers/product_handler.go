package handlers

import (
	"net/http"
	"projectGolang/db"
	"projectGolang/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	var products []models.Product

	// Получаем query-параметры
	categoryID := c.Query("category_id")
	limitParam := c.Query("limit")
	pageParam := c.Query("page")

	// По умолчанию limit и page
	limit := 10
	page := 1

	// Конвертация limit
	if limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil {
			limit = l
		}
	}

	// Конвертация page
	if pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil {
			page = p
		}
	}

	// Построение запроса
	query := db.DB.Limit(limit).Offset((page - 1) * limit)

	// Если есть category_id, добавляем фильтр
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// Выполнение запроса
	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Create(&product)
	c.JSON(http.StatusOK, product)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Save(&product)
	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	db.DB.Delete(&product)
	c.Status(http.StatusNoContent)
}
