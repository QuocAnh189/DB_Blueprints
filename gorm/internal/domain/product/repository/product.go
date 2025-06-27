package repository

import (
	"context"
	"db_blueprints/config"
	db "db_blueprints/gorm/database"
	"db_blueprints/gorm/internal/domain/product/controller/dto"
	"db_blueprints/gorm/internal/model"
	"db_blueprints/pkgs/paging"
)

type IProductRepository interface {
	ListProducts(ctx context.Context, req *dto.ListProductRequest) ([]*model.Product, *paging.Pagination, error)
	GetProductById(ctx context.Context, id int64) (*model.Product, error)
	CreatedProduct(ctx context.Context, product *model.Product) error
	UpdateProduct(ctx context.Context, product *model.Product) error
	DeleteProduct(ctx context.Context, product *model.Product) error
}

type ProductRepository struct {
	db db.IDatabase
}

func NewProductRepository(db db.IDatabase) *ProductRepository {
	return &ProductRepository{db: db}
}

func (pr *ProductRepository) ListProducts(ctx context.Context, req *dto.ListProductRequest) ([]*model.Product, *paging.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	query := make([]db.Query, 0)

	if req.Search != "" {
		query = append(query, db.NewQuery("name ILIKE ?", "%"+req.Search+"%"))
	}

	var order string
	if req.OrderDesc {
		order = "created_at DESC"
	} else {
		order = "created_at"
	}

	if req.OrderBy != "" {
		order = req.OrderBy
		if req.OrderDesc {
			order += " DESC"
		}
	}

	var total int64
	if err := pr.db.Count(ctx, &model.Product{}, &total, db.WithQuery(query...)); err != nil {
		return nil, nil, err
	}

	pagination := paging.NewPagination(req.Page, req.Limit, total)

	var products []*model.Product
	if err := pr.db.Find(
		ctx,
		&products,
		db.WithQuery(query...),
		db.WithLimit(int(pagination.Size)),
		db.WithOffset(int(pagination.Skip)),
		db.WithOrder(order),
		db.WithPreload([]string{"Owner"}),
	); err != nil {
		return nil, nil, err
	}

	return products, pagination, nil
}

func (pr *ProductRepository) GetProductById(ctx context.Context, id int64) (*model.Product, error) {
	var product model.Product
	if err := pr.db.FindById(ctx, id, &product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (pr *ProductRepository) CreatedProduct(ctx context.Context, product *model.Product) error {
	return pr.db.Create(ctx, product)
}

func (pr *ProductRepository) UpdateProduct(ctx context.Context, product *model.Product) error {
	return pr.db.Update(ctx, product)
}

func (pr *ProductRepository) DeleteProduct(ctx context.Context, product *model.Product) error {
	return pr.db.Delete(ctx, product)
}
