package service

import (
	"context"
	"database/sql"
	"fmt"

	"db_blueprints/db_sql/internal/domain/product/controller/dto"
	"db_blueprints/db_sql/internal/domain/product/repository"
	user_repo "db_blueprints/db_sql/internal/domain/user/repository"
	"db_blueprints/db_sql/internal/model"
	"db_blueprints/db_sql/pkgs/paging"
)

type IProductService interface {
	ListProducts(ctx context.Context, req *dto.ListProductRequest) ([]*model.Product, *paging.Pagination, error)
	GetByID(ctx context.Context, id int64) (*model.Product, error)
	CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*model.Product, error)
	UpdateProduct(ctx context.Context, id int64, req *dto.UpdateProductRequest) (*model.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}

type ProductService struct {
	repo      repository.IProductRepository
	user_repo user_repo.IUserRepository
}

func NewProductService(
	repo repository.IProductRepository,
	user_repo user_repo.IUserRepository,
) IProductService {
	return &ProductService{
		repo:      repo,
		user_repo: user_repo,
	}
}

func (s *ProductService) ListProducts(ctx context.Context, req *dto.ListProductRequest) ([]*model.Product, *paging.Pagination, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = paging.DefaultPageSize
	}

	products, total, err := s.repo.List(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("service: failed to list products: %w", err)
	}
	if len(products) == 0 {
		return nil, nil, nil
	}

	// 2. Thu thập các owner_id
	ownerIDs := make([]int64, 0, len(products))
	for _, p := range products {
		ownerIDs = append(ownerIDs, p.OwnerID)
	}

	owners, err := s.user_repo.ListByIDs(ctx, ownerIDs)
	if err != nil {
		return nil, nil, fmt.Errorf("service: failed to get owners for products: %w", err)
	}

	ownerMap := make(map[int64]*model.User, len(owners))
	for _, owner := range owners {
		ownerMap[owner.ID] = owner
	}

	productResponses := make([]*model.Product, 0, len(products))
	for _, p := range products {
		productResp := &model.Product{
			ID:        p.ID,
			Name:      p.Name,
			Price:     p.Price,
			OwnerID:   p.OwnerID,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}

		if owner, ok := ownerMap[p.OwnerID]; ok {
			productResp.Owner = &model.User{
				ID:    owner.ID,
				Name:  owner.Name,
				Email: owner.Email,
			}
		}
		productResponses = append(productResponses, productResp)
	}

	pagination := paging.NewPagination(req.Page, req.Limit, total)
	pagination.TakeAll = req.TakeAll

	return productResponses, pagination, nil
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get product by id: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("service: product with id %d not found", id)
	}

	user, err := s.user_repo.GetByID(ctx, product.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get user by id: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("service: user with id %d not found", id)
	}

	product.Owner = user

	return product, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*model.Product, error) {
	product := &model.Product{
		Name:    req.Name,
		Price:   req.Price,
		OwnerID: req.OwnerID,
	}

	createdProduct, err := s.repo.Create(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("service: failed to create product: %w", err)
	}

	return createdProduct, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int64, req *dto.UpdateProductRequest) (*model.Product, error) {
	productToUpdate, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get product for update: %w", err)
	}
	if productToUpdate == nil {
		return nil, fmt.Errorf("service: product with id %d not found for update", id)
	}

	if req.Name != nil {
		productToUpdate.Name = *req.Name
	}
	if req.Price != nil {
		productToUpdate.Price = *req.Price
	}

	updatedProduct, err := s.repo.Update(ctx, productToUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("service: product with id %d not found on update attempt", id)
		}
		return nil, fmt.Errorf("service: failed to update product: %w", err)
	}

	return updatedProduct, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("service: product with id %d cannot be deleted: %w", id, err)
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("service: product with id %d not found for deletion", id)
		}
		return fmt.Errorf("service: failed to delete product: %w", err)
	}

	return nil
}
