package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"projectGolang/models"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"projectGolang/db"
)

// setupRouter настраивает тестовый Gin router
func setupRouter() *gin.Engine {
	os.Setenv("DB_NAME", "startios_test") // подключение к тестовой базе
	db.InitDB()

	db.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Category{})
	// Очистим таблицу users
	db.DB.Exec("DELETE FROM users")

	r := gin.Default()
	r.POST("/register", Register)
	r.POST("/login", Login)
	return r
}

// TestRegisterSuccess проверяет успешную регистрацию
func TestRegisterSuccess(t *testing.T) {
	router := setupRouter()

	body := `{"name":"Test User","username":"testuser1","password":"123456"}`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Registration successful")
}

// TestLoginSuccess проверяет успешный логин
func TestLoginSuccess(t *testing.T) {
	router := setupRouter()

	// 1. Регистрируем пользователя
	registerBody := `{"name":"Login User","username":"loginuser1","password":"123456"}`
	reqRegister, _ := http.NewRequest("POST", "/register", strings.NewReader(registerBody))
	reqRegister.Header.Set("Content-Type", "application/json")
	wRegister := httptest.NewRecorder()
	router.ServeHTTP(wRegister, reqRegister)
	assert.Equal(t, http.StatusOK, wRegister.Code)

	// 2. Логинимся
	loginBody := `{"username":"loginuser1","password":"123456"}`
	reqLogin, _ := http.NewRequest("POST", "/login", strings.NewReader(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	router.ServeHTTP(wLogin, reqLogin)

	assert.Equal(t, http.StatusOK, wLogin.Code)
	assert.Contains(t, wLogin.Body.String(), "token")
}
