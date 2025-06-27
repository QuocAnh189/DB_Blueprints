package http

import (
	"db_blueprints/db_sql/internal/domain/user/controller/dto"
	"db_blueprints/db_sql/internal/domain/user/service"
	"db_blueprints/db_sql/internal/model"
	"db_blueprints/db_sql/pkgs/response"
	"db_blueprints/db_sql/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.IUserService
}

func NewUserHandler(service service.IUserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	var req dto.ListUserRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Println("Failed to get query", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	users, pagination, err := h.service.ListUsers(c, &req)
	if err != nil {
		log.Println("Failed to get users", err)
		response.Error(c, http.StatusInternalServerError, err, "Failed to get users")
		return
	}

	var res dto.ListUserResponse
	utils.MapStruct(&res.Users, users)
	res.Pagination = pagination

	response.JSON(c, http.StatusOK, res)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	var res model.User

	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("Failed to parse", err)
	}

	user, err := h.service.GetByID(c, userId)
	if err != nil {
		log.Println("Failed to get user", err)
		response.Error(c, http.StatusInternalServerError, err, "Failed to get user")
		return
	}

	utils.MapStruct(&res, user)
	response.JSON(c, http.StatusOK, res)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	user, err := h.service.CreateUser(c, &req)
	if err != nil {
		log.Println("Failed to create user", err)
		response.Error(c, http.StatusInternalServerError, err, "Failed to create user")
		return
	}

	var res dto.CreateUserResponse
	utils.MapStruct(&res.User, user)

	response.JSON(c, http.StatusCreated, res)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req dto.UpdateUserRequest

	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("Failed to parse", err)
	}

	if userId != req.ID {
		response.Error(c, http.StatusBadRequest, nil, "User ID mismatch")
		return
	}

	user, err := h.service.UpdateUser(c, userId, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.UpdateUserResponse
	utils.MapStruct(&res.User, user)

	response.JSON(c, http.StatusOK, res)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("Failed to parse", err)
	}

	err = h.service.DeleteUser(c, userId)

	if err != nil {
		response.Error(c, http.StatusNotFound, err, "Not found")
		return
	}

	response.JSON(c, http.StatusOK, "Delete user successfully")
}
