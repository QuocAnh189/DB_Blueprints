package http

import (
	gorm_db "db_blueprints/blueprints/gorm"
	"db_blueprints/internal/domain/product/repository"
	"db_blueprints/internal/domain/product/service"

	"github.com/gin-gonic/gin"
)

func Routes(
	r *gin.RouterGroup,
	gorm gorm_db.IDatabase,
) {
	productRepository := repository.NewProductRepository(gorm)
	productService := service.NewProductService(productRepository)
	productHandler := NewProductHandler(productService)

	productRoute := r.Group("/products")
	{
		productRoute.GET("", productHandler.GetProducts)
		productRoute.GET("/:id", productHandler.GetProduct)
		productRoute.POST("", productHandler.CreateProduct)
		productRoute.PUT("/:id", productHandler.UpdateProduct)
		productRoute.DELETE("/:id", productHandler.DeleteProduct)
	}
}
