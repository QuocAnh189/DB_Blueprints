package http

import (
	"db_blueprints/internal/domain/user/service"
	"db_blueprints/internal/pkgs/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.IUserService
}

func NewUserHandler(service service.IUserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	response.JSON(c, http.StatusOK, "Test GetUsers")
}

func (h *UserHandler) GetUser(c *gin.Context) {
	response.JSON(c, http.StatusOK, "Test GetUser")
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	response.JSON(c, http.StatusOK, "Test CreateUser")
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	response.JSON(c, http.StatusOK, "Test UpdateUser")
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	response.JSON(c, http.StatusOK, "Test DeleteUser")
}
