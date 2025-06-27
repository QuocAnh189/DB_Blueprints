package service

import (
	"context"
	"database/sql"
	"fmt"

	"db_blueprints/db_sql/internal/domain/user/controller/dto"
	"db_blueprints/db_sql/internal/domain/user/repository"
	"db_blueprints/db_sql/internal/model"
	"db_blueprints/db_sql/pkgs/paging"
)

type IUserService interface {
	ListUsers(ctx context.Context, req *dto.ListUserRequest) ([]*model.User, *paging.Pagination, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error)
	UpdateUser(ctx context.Context, id int64, req *dto.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type UserService struct {
	repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) IUserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) ListUsers(ctx context.Context, req *dto.ListUserRequest) ([]*model.User, *paging.Pagination, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = paging.DefaultPageSize
	}

	users, total, err := s.repo.List(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("service: failed to list users: %w", err)
	}

	pagination := paging.NewPagination(req.Page, req.Limit, total)
	pagination.TakeAll = req.TakeAll

	return users, pagination, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get user by id: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("service: user with id %d not found", id)
	}
	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error) {
	user := &model.User{
		Name:  req.Name,
		Email: req.Email,
	}

	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("service: failed to create user: %w", err)
	}

	return createdUser, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, req *dto.UpdateUserRequest) (*model.User, error) {
	userToUpdate, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get user for update: %w", err)
	}
	if userToUpdate == nil {
		return nil, fmt.Errorf("service: user with id %d not found for update", id)
	}

	if req.Name != nil {
		userToUpdate.Name = *req.Name
	}
	if req.Email != nil {
		userToUpdate.Email = *req.Email
	}

	updatedUser, err := s.repo.Update(ctx, userToUpdate)
	if err != nil {
		return nil, fmt.Errorf("service: failed to update user: %w", err)
	}

	return updatedUser, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("service: user with id %d cannot be deleted: %w", id, err)
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("service: user with id %d not found for deletion", id)
		}
		return fmt.Errorf("service: failed to delete user: %w", err)
	}

	return nil
}
