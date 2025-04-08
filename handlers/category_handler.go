package handlers

import (
	"encoding/json"
	"net/http"
	"projectGolang/models"
)

var categories = []models.Category{
	{ID: 1, Name: "Smartphones"},
	{ID: 2, Name: "Laptops"},
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	json.NewDecoder(r.Body).Decode(&category)
	category.ID = len(categories) + 1
	categories = append(categories, category)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}
