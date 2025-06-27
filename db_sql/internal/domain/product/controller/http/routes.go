package http

import (
	db "db_blueprints/db_sql/database"
	"db_blueprints/db_sql/internal/domain/product/repository"
	"db_blueprints/db_sql/internal/domain/product/service"
	user_repo "db_blueprints/db_sql/internal/domain/user/repository"

	"github.com/gin-gonic/gin"
)

func Routes(
	r *gin.RouterGroup,
	db db.DBTX,
) {
	productRepository := repository.NewProductRepository(db)
	userRepository := user_repo.NewUserRepository(db)
	productService := service.NewProductService(productRepository, userRepository)
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
