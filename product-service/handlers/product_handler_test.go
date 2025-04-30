package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"projectGolang/db"
	userhandlers "projectGolang/handlers"
	"projectGolang/middleware"
	"projectGolang/models"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	os.Setenv("DB_NAME", "startios_test")
	db.InitDB()

	// Миграции
	db.DB.AutoMigrate(&models.Product{}, &models.User{}, &models.Category{})

	// Очистка
	db.DB.Exec("DELETE FROM products")
	db.DB.Exec("DELETE FROM users")

	r := gin.Default()
	r.POST("/register", userhandlers.Register)
	r.POST("/login", userhandlers.Login)

	auth := r.Group("/", middleware.AuthMiddleware())
	auth.GET("/products", GetProducts)
	auth.GET("/products/:id", GetProductByID)
	auth.POST("/products", CreateProduct)
	auth.PUT("/products/:id", UpdateProduct)
	auth.DELETE("/products/:id", DeleteProduct)
	auth.GET("/products/search", SearchProducts)
	auth.GET("/profile", userhandlers.GetProfile)
	auth.PUT("/profile/password", userhandlers.ChangePassword)

	return r
}

func getAuthToken(t *testing.T, router *gin.Engine) string {
	registerBody := `{"name":"Test User","username":"authuser","password":"123456"}`
	reqReg, _ := http.NewRequest("POST", "/register", strings.NewReader(registerBody))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	router.ServeHTTP(wReg, reqReg)

	loginBody := `{"username":"authuser","password":"123456"}`
	reqLogin, _ := http.NewRequest("POST", "/login", strings.NewReader(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	router.ServeHTTP(wLogin, reqLogin)

	var resp map[string]string
	_ = json.Unmarshal(wLogin.Body.Bytes(), &resp)

	return resp["token"]
}

func TestCreateProduct_Success(t *testing.T) {
	router := setupTestRouter()
	token := getAuthToken(t, router)

	productBody := `{"name":"Test Product","category_id":1,"price":100.0}`
	req, _ := http.NewRequest("POST", "/products", strings.NewReader(productBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Product")
}

func TestGetProducts_Success(t *testing.T) {
	router := setupTestRouter()
	token := getAuthToken(t, router)

	productBody := `{"name":"Product 1","category_id":1,"price":99.99}`
	reqCreate, _ := http.NewRequest("POST", "/products", strings.NewReader(productBody))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	wCreate := httptest.NewRecorder()
	router.ServeHTTP(wCreate, reqCreate)

	reqGet, _ := http.NewRequest("GET", "/products", nil)
	reqGet.Header.Set("Authorization", "Bearer "+token)
	wGet := httptest.NewRecorder()
	router.ServeHTTP(wGet, reqGet)

	assert.Equal(t, http.StatusOK, wGet.Code)
	assert.Contains(t, wGet.Body.String(), "Product 1")
}

func TestGetProductByID_Success(t *testing.T) {
	router := setupTestRouter()
	token := getAuthToken(t, router)

	createBody := `{"name":"PlayStation","category_id":1,"price":500}`
	reqCreate, _ := http.NewRequest("POST", "/products", strings.NewReader(createBody))
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	router.ServeHTTP(wCreate, reqCreate)

	var created models.Product
	_ = json.Unmarshal(wCreate.Body.Bytes(), &created)

	reqGet, _ := http.NewRequest("GET", fmt.Sprintf("/products/%d", created.ID), nil)
	reqGet.Header.Set("Authorization", "Bearer "+token)
	wGet := httptest.NewRecorder()
	router.ServeHTTP(wGet, reqGet)

	assert.Equal(t, http.StatusOK, wGet.Code)
	assert.Contains(t, wGet.Body.String(), "PlayStation")
}

func TestUpdateProduct_Success(t *testing.T) {
	router := setupTestRouter()
	token := getAuthToken(t, router)

	createBody := `{"name":"Old Name","category_id":1,"price":10}`
	reqCreate, _ := http.NewRequest("POST", "/products", strings.NewReader(createBody))
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	router.ServeHTTP(wCreate, reqCreate)

	var created models.Product
	_ = json.Unmarshal(wCreate.Body.Bytes(), &created)

	updateBody := `{"name":"New Name","category_id":1,"price":20}`
	reqUpdate, _ := http.NewRequest("PUT", fmt.Sprintf("/products/%d", created.ID), strings.NewReader(updateBody))
	reqUpdate.Header.Set("Authorization", "Bearer "+token)
	reqUpdate.Header.Set("Content-Type", "application/json")
	wUpdate := httptest.NewRecorder()
	router.ServeHTTP(wUpdate, reqUpdate)

	assert.Equal(t, http.StatusOK, wUpdate.Code)
	assert.Contains(t, wUpdate.Body.String(), "New Name")
}

func TestDeleteProduct_Success(t *testing.T) {
	router := setupTestRouter()
	token := getAuthToken(t, router)

	createBody := `{"name":"To Be Deleted","category_id":1,"price":10}`
	reqCreate, _ := http.NewRequest("POST", "/products", strings.NewReader(createBody))
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	router.ServeHTTP(wCreate, reqCreate)

	var created models.Product
	_ = json.Unmarshal(wCreate.Body.Bytes(), &created)

	reqDelete, _ := http.NewRequest("DELETE", fmt.Sprintf("/products/%d", created.ID), nil)
	reqDelete.Header.Set("Authorization", "Bearer "+token)
	wDelete := httptest.NewRecorder()
	router.ServeHTTP(wDelete, reqDelete)

	assert.Equal(t, http.StatusNoContent, wDelete.Code)
}

func TestSearchProduct_Success(t *testing.T) {
	router := setupTestRouter()
	token := getAuthToken(t, router)

	createBody := `{"name":"UniqueProductXYZ","category_id":1,"price":999}`
	reqCreate, _ := http.NewRequest("POST", "/products", strings.NewReader(createBody))
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	router.ServeHTTP(wCreate, reqCreate)

	reqSearch, _ := http.NewRequest("GET", "/products/search?name=xyz", nil)
	reqSearch.Header.Set("Authorization", "Bearer "+token)
	wSearch := httptest.NewRecorder()
	router.ServeHTTP(wSearch, reqSearch)

	assert.Equal(t, http.StatusOK, wSearch.Code)
	assert.Contains(t, wSearch.Body.String(), "UniqueProductXYZ")
}

func TestGetProfile_Success(t *testing.T) {
	router := setupTestRouter()
	token := getAuthToken(t, router)

	req, _ := http.NewRequest("GET", "/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test User")
}

func TestChangePassword_Success(t *testing.T) {
	router := setupTestRouter()
	token := getAuthToken(t, router)

	changeBody := `{
		"old_password": "123456",
		"new_password": "654321"
	}`
	req, _ := http.NewRequest("PUT", "/profile/password", strings.NewReader(changeBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Password changed successfully")
}
