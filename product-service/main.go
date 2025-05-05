package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"projectGolang/db"
	"projectGolang/middleware"
	"projectGolang/product-service/handlers"
)

func main() {
	db.InitDB()
	r := gin.New()
	r.Use(gin.Recovery(), middleware.LoggerMiddleware())

	// 👇 Ашық маршруттар
	r.GET("/products", handlers.GetProducts)
	r.GET("/products/:id", handlers.GetProductByID)
	r.GET("/products/search", handlers.SearchProducts)

	// 👇 Қорғалған маршруттар (токен қажет)
	auth := r.Group("/", middleware.AuthMiddleware())
	auth.GET("/products/profile", handlers.GetProfileFromUserService)
	auth.POST("/products", handlers.CreateProduct)
	auth.PUT("/products/:id", handlers.UpdateProduct)
	auth.DELETE("/products/:id", handlers.DeleteProduct)

	log.Println("✅ ProductService started at :8081")
	r.Run(":8081")
}
