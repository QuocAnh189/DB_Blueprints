package repository

import (
	"context"
	gorm_db "db_blueprints/blueprints/gorm"
	"db_blueprints/internal/config"
	"db_blueprints/internal/domain/product/controller/dto"
	"db_blueprints/internal/model"
	"db_blueprints/internal/pkgs/paging"
)

type IProductGormRepository interface {
	ListProducts(ctx context.Context, req *dto.ListProductRequest) ([]*model.Product, *paging.Pagination, error)
	GetProductById(ctx context.Context, id string) (*model.Product, error)
	CreatedProduct(ctx context.Context, product *model.Product) error
	UpdateProduct(ctx context.Context, product *model.Product) error
	DeleteProduct(ctx context.Context, product *model.Product) error
}

type ProductGormRepository struct {
	gorm_db gorm_db.IDatabase
}

func NewProductRepository(gorm_db gorm_db.IDatabase) *ProductGormRepository {
	return &ProductGormRepository{gorm_db: gorm_db}
}

func (pr *ProductGormRepository) ListProducts(ctx context.Context, req *dto.ListProductRequest) ([]*model.Product, *paging.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	query := make([]gorm_db.Query, 0)

	if req.Search != "" {
		query = append(query, gorm_db.NewQuery("name ILIKE ?", "%"+req.Search+"%"))
	}

	order := "created_at DESC"
	if req.OrderBy != "" {
		order = req.OrderBy
		if req.OrderDesc {
			order += " DESC"
		}
	}

	var total int64
	if err := pr.gorm_db.Count(ctx, &model.Product{}, &total, gorm_db.WithQuery(query...)); err != nil {
		return nil, nil, err
	}

	pagination := paging.NewPagination(req.Page, req.Limit, total)

	var products []*model.Product
	if err := pr.gorm_db.Find(
		ctx,
		&products,
		gorm_db.WithQuery(query...),
		gorm_db.WithLimit(int(pagination.Size)),
		gorm_db.WithOffset(int(pagination.Skip)),
		gorm_db.WithOrder(order),
	); err != nil {
		return nil, nil, err
	}

	return products, pagination, nil
}

func (pr *ProductGormRepository) GetProductById(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product
	if err := pr.gorm_db.FindById(ctx, id, &product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (pr *ProductGormRepository) CreatedProduct(ctx context.Context, product *model.Product) error {
	return pr.gorm_db.Create(ctx, product)
}

func (pr *ProductGormRepository) UpdateProduct(ctx context.Context, product *model.Product) error {
	return pr.gorm_db.Update(ctx, product)
}

func (pr *ProductGormRepository) DeleteProduct(ctx context.Context, product *model.Product) error {
	return pr.gorm_db.Delete(ctx, product)
}
