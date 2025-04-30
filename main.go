package main

import (
	"log"
	"projectGolang/db"
	"projectGolang/handlers"
	"projectGolang/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	r := gin.New()
	r.Use(gin.Recovery(), middleware.LoggerMiddleware())

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	auth := r.Group("/", middleware.AuthMiddleware())

	// Users (protected)
	auth.GET("/users", handlers.GetUsers)
	auth.GET("/users/:id", handlers.GetUserByID)
	auth.POST("/users", handlers.CreateUser)
	auth.PUT("/users/:id", handlers.UpdateUser)
	auth.DELETE("/users/:id", handlers.DeleteUser)
	auth.GET("/profile", handlers.GetProfile)
	auth.PUT("/profile/password", handlers.ChangePassword)

	// Categories (protected)
	auth.GET("/categories", handlers.GetCategories)
	auth.GET("/categories/:id", handlers.GetCategoryByID)
	auth.POST("/categories", handlers.CreateCategory)
	auth.PUT("/categories/:id", handlers.UpdateCategory)
	auth.DELETE("/categories/:id", handlers.DeleteCategory)

	log.Println("Server started at :8080")
	r.Run(":8080")
}
