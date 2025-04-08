package handlers

import (
	"net/http"
	"projectGolang/db"
	"projectGolang/models"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	var categories []models.Category
	db.DB.Find(&categories)
	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Create(&category)
	c.JSON(http.StatusOK, category)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := db.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Save(&category)
	c.JSON(http.StatusOK, category)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := db.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	db.DB.Delete(&category)
	c.Status(http.StatusNoContent)
}
