package repository

import (
	"context"
	"db_blueprints/config"
	db "db_blueprints/gorm/database"
	"db_blueprints/gorm/internal/domain/user/controller/dto"
	"db_blueprints/gorm/internal/model"
	"db_blueprints/pkgs/paging"
)

type IUserRepository interface {
	ListUsers(ctx context.Context, req *dto.ListUserRequest) ([]*model.User, *paging.Pagination, error)
	GetUserById(ctx context.Context, id int64) (*model.User, error)
	CreatedUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, user *model.User) error
}

type UserRepository struct {
	db db.IDatabase
}

func NewUserRepository(db db.IDatabase) *UserRepository {
	return &UserRepository{db: db}
}

func (pr *UserRepository) ListUsers(ctx context.Context, req *dto.ListUserRequest) ([]*model.User, *paging.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	query := make([]db.Query, 0)

	if req.Search != "" {
		query = append(query, db.NewQuery("name LIKE ?", "%"+req.Search+"%"))
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
	if err := pr.db.Count(ctx, &model.User{}, &total, db.WithQuery(query...)); err != nil {
		return nil, nil, err
	}

	pagination := paging.NewPagination(req.Page, req.Limit, total)

	var users []*model.User
	if err := pr.db.Find(
		ctx,
		&users,
		db.WithQuery(query...),
		db.WithLimit(int(pagination.Size)),
		db.WithOffset(int(pagination.Skip)),
		db.WithOrder(order),
	); err != nil {
		return nil, nil, err
	}

	return users, pagination, nil
}

func (pr *UserRepository) GetUserById(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	if err := pr.db.FindById(ctx, id, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (pr *UserRepository) CreatedUser(ctx context.Context, user *model.User) error {
	return pr.db.Create(ctx, user)
}

func (pr *UserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	return pr.db.Update(ctx, user)
}

func (pr *UserRepository) DeleteUser(ctx context.Context, user *model.User) error {
	return pr.db.Delete(ctx, user)
}
