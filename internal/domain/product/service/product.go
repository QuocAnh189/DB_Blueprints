package service

import (
	"context"
	"db_blueprints/internal/domain/product/controller/dto"
	"db_blueprints/internal/domain/product/repository"
	"db_blueprints/internal/model"
	"db_blueprints/internal/pkgs/paging"
	"db_blueprints/internal/utils"
	"log"
)

type IProductService interface {
	ListProducts(ctx context.Context, req *dto.ListProductRequest) ([]*model.Product, *paging.Pagination, error)
	GetProductById(ctx context.Context, id int64) (*model.Product, error)
	CreateProduct(ctx context.Context, req *dto.CreateProductRequest) error
	UpdateProduct(ctx context.Context, req *dto.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, id int64) error
}

type ProductService struct {
	gorm_repo repository.IProductGormRepository
}

func NewProductService(gorm_repo repository.IProductGormRepository) *ProductService {
	return &ProductService{
		gorm_repo: gorm_repo,
	}
}

func (pu *ProductService) ListProducts(ctx context.Context, req *dto.ListProductRequest) ([]*model.Product, *paging.Pagination, error) {
	products, pagination, err := pu.gorm_repo.ListProducts(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return products, pagination, nil
}

func (pu *ProductService) GetProductById(ctx context.Context, id int64) (*model.Product, error) {
	product, err := pu.gorm_repo.GetProductById(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (pu *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) error {
	var product model.Product
	utils.MapStruct(&product, &req)

	err := pu.gorm_repo.CreatedProduct(ctx, &product)
	if err != nil {
		log.Printf("Create fail, error: %s", err)
		return err
	}
	return nil
}

func (pu *ProductService) UpdateProduct(ctx context.Context, req *dto.UpdateProductRequest) error {
	product, err := pu.gorm_repo.GetProductById(ctx, req.ID)
	if err != nil {
		log.Printf("Get fail, error: %s", err)
		return err
	}
	utils.MapStruct(product, req)

	err = pu.gorm_repo.UpdateProduct(ctx, product)
	if err != nil {
		log.Printf("Update fail, id: %d, error: %s", req.ID, err)
		return err
	}

	return nil
}

func (pu *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	product, err := pu.gorm_repo.GetProductById(ctx, id)
	if err != nil {
		return err
	}

	if err := pu.gorm_repo.DeleteProduct(ctx, product); err != nil {
		return err
	}

	return nil
}
