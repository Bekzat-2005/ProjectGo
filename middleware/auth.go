// middleware/auth.go
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

var JwtKey = []byte("secret-key")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if uidFloat, ok := claims["user_id"].(float64); ok {
				c.Set("user_id", int(uidFloat)) // ✅ Надёжное преобразование
			}
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
		}
	}
}

//func AuthMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		authHeader := c.GetHeader("Authorization")
//		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
//			c.Abort()
//			return
//		}
//
//		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
//
//		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//			return JwtKey, nil
//		})
//
//		if err != nil || !token.Valid {
//			c.JSON(http.StatusUnauthorized, gin.H{
//				"error":   "Invalid token",
//				"details": err.Error(),
//			})
//			c.Abort()
//			return
//		}
//
//		// Токен жарамды — енді claims ішінен user_id аламыз
//		if claims, ok := token.Claims.(jwt.MapClaims); ok {
//			if uidFloat, ok := claims["user_id"].(float64); ok {
//				c.Set("user_id", uint(uidFloat)) // ✅ user_id-ны uint түрінде сақтаймыз
//			} else {
//				c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
//				c.Abort()
//				return
//			}
//		}
//
//		c.Next()
//	}
//}
