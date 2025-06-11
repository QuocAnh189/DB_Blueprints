package http

import (
	"db_blueprints/internal/domain/product/service"
	"db_blueprints/internal/pkgs/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service service.IProductService
}

func NewProductHandler(service service.IProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	response.JSON(c, http.StatusOK, "Test GetProducts")
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	response.JSON(c, http.StatusOK, "Test GetProduct")
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	response.JSON(c, http.StatusOK, "Test CreateProduct")
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	response.JSON(c, http.StatusOK, "Test UpdateProduct")
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	response.JSON(c, http.StatusOK, "Test DeleteProduct")
}
