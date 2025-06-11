package repository

import (
	gorm_db "db_blueprints/blueprints/gorm"
)

type IProductGormRepository interface {
}

type ProductGormRepository struct {
	gorm_db gorm_db.IDatabase
}

func NewProductRepository(gorm_db gorm_db.IDatabase) *ProductGormRepository {
	return &ProductGormRepository{gorm_db: gorm_db}
}
