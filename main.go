package main

import (
	"log"
	"projectGolang/db"
	"projectGolang/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	r := gin.Default()

	// Users
	r.GET("/users", handlers.GetUsers)
	r.GET("/users/:id", handlers.GetUserByID)
	r.POST("/users", handlers.CreateUser)
	r.PUT("/users/:id", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)

	// Categories
	r.GET("/categories", handlers.GetCategories)
	r.GET("/categories/:id", handlers.GetCategoryByID)
	r.POST("/categories", handlers.CreateCategory)
	r.PUT("/categories/:id", handlers.UpdateCategory)
	r.DELETE("/categories/:id", handlers.DeleteCategory)

	// Products
	r.GET("/products", handlers.GetProducts)
	r.POST("/products", handlers.CreateProduct)
	r.GET("/products/:id", handlers.GetProductByID)
	r.PUT("/products/:id", handlers.UpdateProduct)
	r.DELETE("/products/:id", handlers.DeleteProduct)

	log.Println("Server started at :8080")
	r.Run(":8080")
}
