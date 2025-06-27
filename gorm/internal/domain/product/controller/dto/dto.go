package dto

import (
	"db_blueprints/gorm/internal/domain/user/controller/dto"
	"db_blueprints/pkgs/paging"
)

type Product struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Price     float64  `json:"price"`
	OwnerID   int64    `json:"owner_id"`
	Owner     dto.User `json:"user"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type ListProductRequest struct {
	Search    string `json:"search,omitempty" form:"search"`
	Page      int64  `json:"-" form:"page"`
	Limit     int64  `json:"-" form:"size"`
	OrderBy   string `json:"-" form:"order_by"`
	OrderDesc bool   `json:"-" form:"order_desc"`
	TakeAll   bool   `json:"-" form:"take_all"`
}

type ListProductResponse struct {
	Products   []*Product         `json:"items"`
	Pagination *paging.Pagination `json:"metadata"`
}

type CreateProductRequest struct {
	OwnerID int64   `json:"owner_id"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
}

type CreateProductResponse struct {
	Product *Product `json:"product"`
}

type UpdateProductRequest struct {
	ID    int64    `json:"id"`
	Owner *int64   `json:"owner_id"`
	Name  *string  `json:"name"`
	Price *float64 `json:"price"`
}

type UpdateProductResponse struct {
	Product *Product `json:"product"`
}
