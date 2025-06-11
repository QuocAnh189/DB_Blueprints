package repository

import (
	gorm_db "db_blueprints/blueprints/gorm"
)

type IUserGormRepository interface {
}

type UserGormRepository struct {
	gorm_db gorm_db.IDatabase
}

func NewUserRepository(gorm_db gorm_db.IDatabase) *UserGormRepository {
	return &UserGormRepository{gorm_db: gorm_db}
}
