package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"projectGolang/product-service/db"
	"projectGolang/product-service/handlers"
	"projectGolang/product-service/middleware"
)

func main() {
	db.InitDB() // ✅ Не забудь вызвать инициализацию базы

	r := gin.New()
	r.Use(gin.Recovery(), middleware.LoggerMiddleware())

	r.GET("/products", handlers.GetProducts)
	r.GET("/products/:id", handlers.GetProductByID)
	r.POST("/products", handlers.CreateProduct)
	r.PUT("/products/:id", handlers.UpdateProduct)
	r.DELETE("/products/:id", handlers.DeleteProduct)
	r.GET("/products/search", handlers.SearchProducts)

	r.GET("/products/profile", handlers.GetProfileFromUserService)

	log.Println("✅ ProductService started at :8081")
	r.Run(":8081")
}
