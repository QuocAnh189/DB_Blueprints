package service

import (
	"context"
	"db_blueprints/gorm/internal/domain/user/controller/dto"
	"db_blueprints/gorm/internal/domain/user/repository"
	"db_blueprints/gorm/internal/model"
	"db_blueprints/gorm/utils"
	"db_blueprints/pkgs/paging"
	"log"
)

type IUserService interface {
	ListUsers(ctx context.Context, req *dto.ListUserRequest) ([]*model.User, *paging.Pagination, error)
	GetUserById(ctx context.Context, id int64) (*model.User, error)
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error)
	UpdateUser(ctx context.Context, req *dto.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type UserService struct {
	repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (pu *UserService) ListUsers(ctx context.Context, req *dto.ListUserRequest) ([]*model.User, *paging.Pagination, error) {
	users, pagination, err := pu.repo.ListUsers(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return users, pagination, nil
}

func (pu *UserService) GetUserById(ctx context.Context, id int64) (*model.User, error) {
	User, err := pu.repo.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return User, nil
}

func (pu *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error) {
	var user model.User
	utils.MapStruct(&user, &req)

	err := pu.repo.CreatedUser(ctx, &user)
	if err != nil {
		log.Printf("Create fail, error: %s", err)
		return nil, err
	}
	return &user, nil
}

func (pu *UserService) UpdateUser(ctx context.Context, req *dto.UpdateUserRequest) (*model.User, error) {
	user, err := pu.repo.GetUserById(ctx, req.ID)
	if err != nil {
		log.Printf("Get fail, error: %s", err)
		return nil, err
	}
	utils.MapStruct(user, req)

	err = pu.repo.UpdateUser(ctx, user)
	if err != nil {
		log.Printf("Update fail, id: %d, error: %s", req.ID, err)
		return nil, err
	}

	return user, nil
}

func (pu *UserService) DeleteUser(ctx context.Context, id int64) error {
	User, err := pu.repo.GetUserById(ctx, id)
	if err != nil {
		return err
	}

	if err := pu.repo.DeleteUser(ctx, User); err != nil {
		return err
	}

	return nil
}
