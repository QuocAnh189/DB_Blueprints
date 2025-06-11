package service

import (
	"db_blueprints/internal/domain/product/repository"
)

type IProductService interface{}

type ProductService struct {
	gorm_repo repository.IProductGormRepository
}

func NewProductService(gorm_repo repository.IProductGormRepository) *ProductService {
	return &ProductService{
		gorm_repo: gorm_repo,
	}
}
