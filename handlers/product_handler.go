package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"projectGolang/models"
)

var products = []models.Product{
	{ID: 1, Name: "iPhone 14", CategoryID: 1, Price: 999.99},
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Читаем query параметры
	categoryIDStr := r.URL.Query().Get("category_id")
	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")

	var filteredProducts []models.Product

	// Фильтрация по category_id
	if categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err == nil {
			for _, product := range products {
				if product.CategoryID == categoryID {
					filteredProducts = append(filteredProducts, product)
				}
			}
		}
	} else {
		// Если фильтра нет, берем все продукты
		filteredProducts = products
	}

	// Пагинация
	limit := len(filteredProducts) // по умолчанию все
	page := 1                      // по умолчанию первая страница

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	start := (page - 1) * limit
	end := start + limit

	// Проверяем границы массива
	if start > len(filteredProducts) {
		start = len(filteredProducts)
	}
	if end > len(filteredProducts) {
		end = len(filteredProducts)
	}

	paginatedProducts := filteredProducts[start:end]

	json.NewEncoder(w).Encode(paginatedProducts)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for _, product := range products {
		if product.ID == id {
			json.NewEncoder(w).Encode(product)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	json.NewDecoder(r.Body).Decode(&product)
	product.ID = len(products) + 1
	products = append(products, product)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i, product := range products {
		if product.ID == id {
			json.NewDecoder(r.Body).Decode(&products[i])
			products[i].ID = id
			json.NewEncoder(w).Encode(products[i])
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i, product := range products {
		if product.ID == id {
			products = append(products[:i], products[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
