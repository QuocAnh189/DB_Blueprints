package service

import (
	"db_blueprints/internal/domain/user/repository"
)

type IUserService interface {
}

type UserService struct {
	gorm_repo repository.IUserGormRepository
}

func NewUserService(gorm_repo repository.IUserGormRepository) *UserService {
	return &UserService{
		gorm_repo: gorm_repo,
	}
}
