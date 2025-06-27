package service

import (
	"context"
	"db_blueprints/gorm/internal/domain/product/controller/dto"
	"db_blueprints/gorm/internal/domain/product/repository"
	user_repo "db_blueprints/gorm/internal/domain/user/repository"
	"db_blueprints/gorm/internal/model"
	"db_blueprints/gorm/pkgs/paging"
	"db_blueprints/gorm/utils"
	"log"
)

type IProductService interface {
	ListProducts(ctx context.Context, req *dto.ListProductRequest) ([]*model.Product, *paging.Pagination, error)
	GetProductById(ctx context.Context, id int64) (*model.Product, error)
	CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*model.Product, error)
	UpdateProduct(ctx context.Context, req *dto.UpdateProductRequest) (*model.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}

type ProductService struct {
	repo      repository.IProductRepository
	user_repo user_repo.IUserRepository
}

func NewProductService(
	repo repository.IProductRepository,
	user_repo user_repo.IUserRepository,
) *ProductService {
	return &ProductService{
		repo:      repo,
		user_repo: user_repo,
	}
}

func (pu *ProductService) ListProducts(ctx context.Context, req *dto.ListProductRequest) ([]*model.Product, *paging.Pagination, error) {
	products, pagination, err := pu.repo.ListProducts(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return products, pagination, nil
}

func (pu *ProductService) GetProductById(ctx context.Context, id int64) (*model.Product, error) {
	product, err := pu.repo.GetProductById(ctx, id)
	if err != nil {
		return nil, err
	}

	println("user id", product.OwnerID)
	user, err := pu.user_repo.GetUserById(ctx, product.OwnerID)
	if err != nil {
		log.Printf("Get user by id %d fail, error: %s", product.OwnerID, err)
		return nil, err
	}

	if user != nil {
		product.Owner = *user
	}

	return product, nil
}

func (pu *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*model.Product, error) {
	var product model.Product
	utils.MapStruct(&product, &req)

	err := pu.repo.CreatedProduct(ctx, &product)
	if err != nil {
		log.Printf("Create fail, error: %s", err)
		return nil, err
	}
	return &product, nil
}

func (pu *ProductService) UpdateProduct(ctx context.Context, req *dto.UpdateProductRequest) (*model.Product, error) {
	product, err := pu.repo.GetProductById(ctx, req.ID)
	if err != nil {
		log.Printf("Get fail, error: %s", err)
		return nil, err
	}
	utils.MapStruct(product, req)

	err = pu.repo.UpdateProduct(ctx, product)
	if err != nil {
		log.Printf("Update fail, id: %d, error: %s", req.ID, err)
		return nil, err
	}

	return product, nil
}

func (pu *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	product, err := pu.repo.GetProductById(ctx, id)
	if err != nil {
		return err
	}

	if err := pu.repo.DeleteProduct(ctx, product); err != nil {
		return err
	}

	return nil
}
