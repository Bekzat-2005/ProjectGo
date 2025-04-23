package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"projectGolang/db"
)

func setupRouter() *gin.Engine {
	db.InitDB()
	db.DB.Exec("DELETE FROM users") // очистим пользователей

	r := gin.Default()
	r.POST("/register", Register)
	return r
}

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
func TestLoginSuccess(t *testing.T) {
	router := gin.Default()
	db.InitDB()
	db.DB.Exec("DELETE FROM users")

	router.POST("/register", Register)
	router.POST("/login", Login)

	// 1. Сначала регистрируем пользователя
	registerBody := `{"name":"Login User","username":"loginuser1","password":"123456"}`
	reqRegister, _ := http.NewRequest("POST", "/register", strings.NewReader(registerBody))
	reqRegister.Header.Set("Content-Type", "application/json")
	wRegister := httptest.NewRecorder()
	router.ServeHTTP(wRegister, reqRegister)
	assert.Equal(t, http.StatusOK, wRegister.Code)

	// 2. Затем логинимся
	loginBody := `{"username":"loginuser1","password":"123456"}`
	reqLogin, _ := http.NewRequest("POST", "/login", strings.NewReader(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	router.ServeHTTP(wLogin, reqLogin)

	assert.Equal(t, http.StatusOK, wLogin.Code)
	assert.Contains(t, wLogin.Body.String(), "token")
}
