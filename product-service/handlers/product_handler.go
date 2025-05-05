package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"projectGolang/db"
	"projectGolang/models"
	"projectGolang/product-service/client"
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

	// 📦 Получаем токен
	authHeader := c.GetHeader("Authorization")
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		token = authHeader
	}

	// ⬅️ Формируем ответ с user-информацией
	var response []gin.H
	for _, p := range products {
		fmt.Println("🔎 Requesting user with ID:", p.UserID)
		userInfo, err := client.GetUserByID(p.UserID, token)
		if err != nil {
			userInfo = map[string]interface{}{"error": "Failed to fetch user"}
		}

		response = append(response, gin.H{
			"id":          p.ID,
			"name":        p.Name,
			"category_id": p.CategoryID,
			"price":       p.Price,
			"user":        userInfo,
		})
	}

	c.JSON(http.StatusOK, response)
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Authorization header-ден токенді бөліп алу
	authHeader := c.GetHeader("Authorization")
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		token = authHeader
	}

	// 👇 user-service-тен user алу
	userData, err := client.GetUserByID(product.UserID, token)
	if err != nil {
		userData = map[string]interface{}{"error": "User fetch failed"}
	}

	// 👇 userData-ны JSON строкаға айналдыру
	userJSON, err := json.Marshal(userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize user"})
		return
	}

	// ✅ жауап
	c.JSON(http.StatusOK, gin.H{
		"id":          product.ID,
		"name":        product.Name,
		"category_id": product.CategoryID,
		"price":       product.Price,
		"user":        string(userJSON), // 🔥 осы жерде строка
	})
}

func CreateProduct(c *gin.Context) {
	var product models.Product

	// 1. JSON-нан product мәліметтерін оқу
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. user_id-ны токеннен алу и преобразование
	if userID, ok := c.Get("user_id"); ok {
		// Преобразуем в int, а затем в uint
		if uid, ok := userID.(int); ok {
			product.UserID = uint(uid)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type in token"})
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	// 3. Базаға сақтау
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
func SearchProducts(c *gin.Context) {
	name := c.Query("name")

	var products []models.Product
	query := db.DB

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	query.Find(&products)

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// 👇 добавь этот хендлер к остальным
func GetProfileFromUserService(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// Обрезаем "Bearer " если есть
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		token = authHeader
	}

	profile, err := client.GetUserProfile(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}
