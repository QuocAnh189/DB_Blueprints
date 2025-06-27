package http

import (
	"log"
	"net/http"
	"strconv"

	"db_blueprints/db_sql/internal/domain/product/controller/dto"
	"db_blueprints/db_sql/internal/domain/product/service"
	"db_blueprints/db_sql/internal/model"
	"db_blueprints/db_sql/pkgs/response"
	"db_blueprints/db_sql/utils"

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
		log.Printf("Failed to bind query parameters: %v", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	products, pagination, err := h.service.ListProducts(c, &req)
	if err != nil {
		log.Printf("Failed to get products: %v", err)
		response.Error(c, http.StatusInternalServerError, err, "Failed to get products")
		return
	}

	var res dto.ListProductResponse
	utils.MapStruct(&res.Products, products)
	res.Pagination = pagination

	response.JSON(c, http.StatusOK, res)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	var res model.Product

	productId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Failed to parse product ID from path: %v", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid product ID")
		return
	}

	product, err := h.service.GetByID(c, productId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err, "Failed to get product")
		return
	}

	utils.MapStruct(&res, product)
	response.JSON(c, http.StatusOK, res)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	product, err := h.service.CreateProduct(c, &req)
	if err != nil {
		log.Printf("Failed to create product: %v", err)
		response.Error(c, http.StatusInternalServerError, err, "Failed to create product")
		return
	}

	var res dto.CreateProductResponse
	utils.MapStruct(&res.Product, product)

	response.JSON(c, http.StatusCreated, res)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var req dto.UpdateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	productId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Failed to parse product ID from path: %v", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid product ID")
		return
	}

	product, err := h.service.UpdateProduct(c, productId, &req)
	if err != nil {
		log.Printf("Failed to update product: %v", err)
		response.Error(c, http.StatusInternalServerError, err, "Failed to update product")
		return
	}

	var res dto.UpdateProductResponse
	utils.MapStruct(&res.Product, product)

	response.JSON(c, http.StatusOK, res)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Failed to parse product ID from path: %v", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid product ID")
		return
	}

	err = h.service.DeleteProduct(c, productId)
	if err != nil {
		log.Printf("Failed to delete product: %v", err)
		response.Error(c, http.StatusNotFound, err, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, gin.H{"message": "Delete product successfully"})
}
