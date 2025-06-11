package http

import (
	gorm_db "db_blueprints/blueprints/gorm"
	"db_blueprints/internal/domain/user/repository"
	"db_blueprints/internal/domain/user/service"

	"github.com/gin-gonic/gin"
)

func Routes(
	r *gin.RouterGroup,
	gorm gorm_db.IDatabase,
) {
	userRepository := repository.NewUserRepository(gorm)
	userService := service.NewUserService(userRepository)
	userHandler := NewUserHandler(userService)

	userRoute := r.Group("/users")
	{
		userRoute.GET("", userHandler.GetUsers)
		userRoute.GET("/:id", userHandler.GetUser)
		userRoute.POST("", userHandler.CreateUser)
		userRoute.PUT("/:id", userHandler.UpdateUser)
		userRoute.DELETE("/:id", userHandler.DeleteUser)
	}
}
