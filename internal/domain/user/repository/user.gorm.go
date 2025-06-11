package repository

import (
	"context"
	gorm_db "db_blueprints/blueprints/gorm"
	"db_blueprints/internal/config"
	"db_blueprints/internal/domain/user/controller/dto"
	"db_blueprints/internal/model"
	"db_blueprints/internal/pkgs/paging"
)

type IUserGormRepository interface {
	ListUsers(ctx context.Context, req *dto.ListUserRequest) ([]*model.User, *paging.Pagination, error)
	GetUserById(ctx context.Context, id int64) (*model.User, error)
	CreatedUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, user *model.User) error
}

type UserGormRepository struct {
	gorm_db gorm_db.IDatabase
}

func NewUserRepository(gorm_db gorm_db.IDatabase) *UserGormRepository {
	return &UserGormRepository{gorm_db: gorm_db}
}

func (pr *UserGormRepository) ListUsers(ctx context.Context, req *dto.ListUserRequest) ([]*model.User, *paging.Pagination, error) {
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
	if err := pr.gorm_db.Count(ctx, &model.User{}, &total, gorm_db.WithQuery(query...)); err != nil {
		return nil, nil, err
	}

	pagination := paging.NewPagination(req.Page, req.Limit, total)

	var users []*model.User
	if err := pr.gorm_db.Find(
		ctx,
		&users,
		gorm_db.WithQuery(query...),
		gorm_db.WithLimit(int(pagination.Size)),
		gorm_db.WithOffset(int(pagination.Skip)),
		gorm_db.WithOrder(order),
	); err != nil {
		return nil, nil, err
	}

	return users, pagination, nil
}

func (pr *UserGormRepository) GetUserById(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	if err := pr.gorm_db.FindById(ctx, id, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (pr *UserGormRepository) CreatedUser(ctx context.Context, user *model.User) error {
	return pr.gorm_db.Create(ctx, user)
}

func (pr *UserGormRepository) UpdateUser(ctx context.Context, user *model.User) error {
	return pr.gorm_db.Update(ctx, user)
}

func (pr *UserGormRepository) DeleteUser(ctx context.Context, user *model.User) error {
	return pr.gorm_db.Delete(ctx, user)
}
