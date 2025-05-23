package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"projectGolang/db"
	"projectGolang/handlers" // ← сюда входят Register и Login
	"projectGolang/middleware"
	"projectGolang/models"
	productHandlers "projectGolang/product-service/handlers"
)

func extractToken(body string) string {
	var result map[string]string
	_ = json.Unmarshal([]byte(body), &result)
	return result["token"]
}

func setupProductTestRouter() *gin.Engine {
	os.Setenv("DB_NAME", "startios_test")
	db.InitDB()
	db.DB.Exec("DELETE FROM users")
	db.DB.Exec("DELETE FROM categories")
	db.DB.Exec("DELETE FROM products")
	db.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Category{})

	r := gin.Default()
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	auth := r.Group("/", middleware.AuthMiddleware())
	auth.POST("/products", productHandlers.CreateProduct)
	auth.GET("/products", productHandlers.GetProducts)
	auth.GET("/products/:id", productHandlers.GetProductByID)
	auth.PUT("/products/:id", productHandlers.UpdateProduct)
	auth.DELETE("/products/:id", productHandlers.DeleteProduct)
	auth.GET("/products/search", productHandlers.SearchProducts)

	return r
}

func createUserAndLogin(t *testing.T, router *gin.Engine) string {
	rb := `{"name":"Test User","username":"testuser","password":"123456"}`
	rr, _ := http.NewRequest("POST", "/register", strings.NewReader(rb))
	rr.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, rr)

	lb := `{"username":"testuser","password":"123456"}`
	rl, _ := http.NewRequest("POST", "/login", strings.NewReader(lb))
	rl.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, rl)
	return extractToken(w2.Body.String())
}

func TestCreateProductSuccess(t *testing.T) {
	router := setupProductTestRouter()
	token := createUserAndLogin(t, router)

	cat := models.Category{Name: "Electronics"}
	db.DB.Create(&cat)
	body := fmt.Sprintf(`{"name":"Phone","category_id":%d,"price":599.99}`, cat.ID)
	r, _ := http.NewRequest("POST", "/products", strings.NewReader(body))
	r.Header.Set("Authorization", "Bearer "+token)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateProductWithoutToken(t *testing.T) {
	router := setupProductTestRouter()
	body := `{"name":"NoAuth","category_id":1,"price":100}`
	r, _ := http.NewRequest("POST", "/products", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetProductsSuccess(t *testing.T) {
	router := setupProductTestRouter()
	token := createUserAndLogin(t, router)
	r, _ := http.NewRequest("GET", "/products", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetProductByIDNotFound(t *testing.T) {
	router := setupProductTestRouter()
	token := createUserAndLogin(t, router)
	r, _ := http.NewRequest("GET", "/products/999", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSearchProductsByName(t *testing.T) {
	router := setupProductTestRouter()
	token := createUserAndLogin(t, router)

	cat := models.Category{Name: "Books"}
	db.DB.Create(&cat)
	db.DB.Create(&models.Product{Name: "Go Book", CategoryID: cat.ID, Price: 49.99, UserID: 1})

	r, _ := http.NewRequest("GET", "/products/search?name=Go", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Go Book")
}

func TestUpdateProductSuccess(t *testing.T) {
	router := setupProductTestRouter()
	token := createUserAndLogin(t, router)

	cat := models.Category{Name: "Tech"}
	db.DB.Create(&cat)
	product := models.Product{Name: "Old Name", CategoryID: cat.ID, Price: 100, UserID: 1}
	db.DB.Create(&product)

	body := `{"name":"Updated Name","category_id":1,"price":150}`
	r, _ := http.NewRequest("PUT", fmt.Sprintf("/products/%d", product.ID), strings.NewReader(body))
	r.Header.Set("Authorization", "Bearer "+token)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteProductSuccess(t *testing.T) {
	router := setupProductTestRouter()
	token := createUserAndLogin(t, router)

	cat := models.Category{Name: "Delete"}
	db.DB.Create(&cat)
	product := models.Product{Name: "Delete Me", CategoryID: cat.ID, Price: 10, UserID: 1}
	db.DB.Create(&product)

	r, _ := http.NewRequest("DELETE", fmt.Sprintf("/products/%d", product.ID), nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
func TestCreateProductInvalidJSON(t *testing.T) {
	router := setupProductTestRouter()
	token := createUserAndLogin(t, router)

	body := `{"name":123, "category_id":"wrong", "price":"text"}`
	r, _ := http.NewRequest("POST", "/products", strings.NewReader(body))
	r.Header.Set("Authorization", "Bearer "+token)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
