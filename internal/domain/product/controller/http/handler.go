package http

import (
	"db_blueprints/internal/domain/product/controller/dto"
	"db_blueprints/internal/domain/product/service"
	"db_blueprints/internal/model"
	"db_blueprints/internal/pkgs/response"
	"db_blueprints/internal/utils"
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
	var res model.Product

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

	if err := h.service.CreateProduct(c, &req); err != nil {
		log.Println("Failed to create product", err)
		response.Error(c, http.StatusInternalServerError, err, "Failed to create product")
		return
	}

	response.JSON(c, http.StatusCreated, "Create product successfully")
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

	if err := h.service.UpdateProduct(c, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	response.JSON(c, http.StatusOK, "Update product successfully")
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

	response.JSON(c, http.StatusOK, "Delete products successfully")
}
