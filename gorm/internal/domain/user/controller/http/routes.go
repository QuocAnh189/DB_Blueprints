package http

import (
	db "db_blueprints/gorm/database"
	"db_blueprints/gorm/internal/domain/user/repository"
	"db_blueprints/gorm/internal/domain/user/service"

	"github.com/gin-gonic/gin"
)

func Routes(
	r *gin.RouterGroup,
	db db.IDatabase,
) {
	userRepository := repository.NewUserRepository(db)
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
