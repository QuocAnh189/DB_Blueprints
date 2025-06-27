package http

import (
	db "db_blueprints/db_sql/database"
	"db_blueprints/db_sql/internal/domain/user/repository"
	"db_blueprints/db_sql/internal/domain/user/service"

	"github.com/gin-gonic/gin"
)

func Routes(
	r *gin.RouterGroup,
	db db.DBTX,
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
