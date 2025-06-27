package http

import (
	"db_blueprints/gorm/internal/domain/product/controller/dto"
	"db_blueprints/gorm/internal/domain/product/service"
	"db_blueprints/gorm/utils"
	"db_blueprints/pkgs/response"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service service.IProductService
}

func NewProductHandler(service service.IProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	var req dto.ListProductRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Println("Failed to get query", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	products, pagination, err := h.service.ListProducts(c, &req)
	if err != nil {
		log.Println("Failed to get products", err)
		response.Error(c, http.StatusInternalServerError, err, "Failed to get products")
		return
	}

	var res dto.ListProductResponse
	utils.MapStruct(&res.Products, products)
	res.Pagination = pagination

	response.JSON(c, http.StatusOK, res)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	var res dto.Product

	productId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("Failed to parse", err)
	}

	product, err := h.service.GetProductById(c, productId)
	if err != nil {
		log.Println("Failed to get product", err)
		response.Error(c, http.StatusInternalServerError, err, "Failed to get product")
		return
	}

	utils.MapStruct(&res, product)
	response.JSON(c, http.StatusOK, res)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	product, err := h.service.CreateProduct(c, &req)
	if err != nil {
		log.Println("Failed to create product", err)
		response.Error(c, http.StatusInternalServerError, err, "Failed to create product")
		return
	}

	var res dto.CreateProductResponse
	utils.MapStruct(&res.Product, product)

	response.JSON(c, http.StatusCreated, res)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var req dto.UpdateProductRequest

	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	productId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("Failed to parse", err)
	}

	if productId != req.ID {
		response.Error(c, http.StatusBadRequest, nil, "Product ID mismatch")
		return
	}

	product, err := h.service.UpdateProduct(c, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.UpdateProductResponse
	utils.MapStruct(&res.Product, product)

	response.JSON(c, http.StatusOK, res)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("Failed to parse", err)
	}

	err = h.service.DeleteProduct(c, productId)

	if err != nil {
		response.Error(c, http.StatusNotFound, err, "Not found")
		return
	}

	response.JSON(c, http.StatusOK, "Delete user successfully")
}
