package repository

import (
	"context"
	"database/sql"
	"db_blueprints/db_sql/database"
	"db_blueprints/db_sql/internal/domain/product/controller/dto"
	"db_blueprints/db_sql/internal/model"
	"fmt"
	"strings"
)

type IProductRepository interface {
	GetByID(ctx context.Context, id int64) (*model.Product, error)
	Create(ctx context.Context, product *model.Product) (*model.Product, error)
	Update(ctx context.Context, product *model.Product) (*model.Product, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, req *dto.ListProductRequest) ([]*model.Product, int64, error)
}

type ProductRepository struct {
	db database.DBTX
}

func NewProductRepository(db database.DBTX) IProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	query := "SELECT id, name, price, owner_id, created_at, updated_at FROM products WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)
	var p model.Product
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.OwnerID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get product by id: %w", err)
	}
	return &p, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	query := "INSERT INTO products (name, price, owner_id, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())"
	result, err := r.db.ExecContext(ctx, query, product.Name, product.Price, product.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get last insert id for product: %w", err)
	}
	product.ID = id
	return product, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *model.Product) (*model.Product, error) {
	query := "UPDATE products SET name = ?, price = ?, updated_at = NOW() WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, product.Name, product.Price, product.ID)

	if err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("get rows affected for product update: %w", err)
	}

	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return product, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM products WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected for product delete: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ProductRepository) List(ctx context.Context, req *dto.ListProductRequest) ([]*model.Product, int64, error) {
	var total int64
	countQueryBuilder := strings.Builder{}
	countQueryBuilder.WriteString("SELECT COUNT(id) FROM products WHERE 1=1")
	args := []interface{}{}

	if req.Search != "" {
		countQueryBuilder.WriteString(" AND name LIKE ?")
		searchPattern := "%" + req.Search + "%"
		args = append(args, searchPattern)
	}
	err := r.db.QueryRowContext(ctx, countQueryBuilder.String(), args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}
	if total == 0 {
		return []*model.Product{}, 0, nil
	}

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString("SELECT id, name, price, owner_id, created_at, updated_at FROM products WHERE 1=1")
	if req.Search != "" {
		queryBuilder.WriteString(" AND name LIKE ?")
	}

	orderBy := "id"
	allowedOrderBys := map[string]string{
		"name":       "name",
		"price":      "price",
		"created_at": "created_at",
	}
	if col, ok := allowedOrderBys[req.OrderBy]; ok {
		orderBy = col
	}
	queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s", orderBy))
	if req.OrderDesc {
		queryBuilder.WriteString(" DESC")
	} else {
		queryBuilder.WriteString(" ASC")
	}

	if !req.TakeAll {
		queryBuilder.WriteString(" LIMIT ? OFFSET ?")
		offset := (req.Page - 1) * req.Limit
		args = append(args, req.Limit, offset)
	}

	rows, err := r.db.QueryContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}

	scanProduct := func(rows *sql.Rows) (*model.Product, error) {
		var p model.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.OwnerID, &p.CreatedAt, &p.UpdatedAt)
		return &p, err
	}

	products, err := database.ScanRows(rows, scanProduct)
	if err != nil {
		return nil, 0, fmt.Errorf("scan products: %w", err)
	}

	return products, total, nil
}
