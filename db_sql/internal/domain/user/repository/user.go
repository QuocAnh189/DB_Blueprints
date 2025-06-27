package repository

import (
	"context"
	"database/sql"
	"db_blueprints/db_sql/database"
	"db_blueprints/db_sql/internal/domain/user/controller/dto"
	"db_blueprints/db_sql/internal/model"
	"fmt"
	"strings"
)

type IUserRepository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, req *dto.ListUserRequest) ([]*model.User, int64, error)
	ListByIDs(ctx context.Context, ids []int64) ([]*model.User, error)
}

type UserRepository struct {
	db database.DBTX
}

func NewUserRepository(db database.DBTX) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)
	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	query := "INSERT INTO users (name, email, created_at, updated_at) VALUES (?, ?, NOW(), NOW())"
	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get last insert id: %w", err)
	}
	user.ID = id
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	query := "UPDATE users SET name = ?, email = ?, updated_at = NOW() WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.ID)

	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context, req *dto.ListUserRequest) ([]*model.User, int64, error) {
	var total int64
	countQueryBuilder := strings.Builder{}
	countQueryBuilder.WriteString("SELECT COUNT(id) FROM users WHERE 1=1")
	args := []interface{}{}
	if req.Search != "" {
		countQueryBuilder.WriteString(" AND (name LIKE ? OR email LIKE ?)")
		searchPattern := "%" + req.Search + "%"
		args = append(args, searchPattern, searchPattern)
	}
	err := r.db.QueryRowContext(ctx, countQueryBuilder.String(), args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	if total == 0 {
		return []*model.User{}, 0, nil
	}

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString("SELECT id, name, email, created_at, updated_at FROM users WHERE 1=1")
	if req.Search != "" {
		queryBuilder.WriteString(" AND (name LIKE ? OR email LIKE ?)")
	}

	orderBy := "id"
	allowedOrderBys := map[string]string{"name": "name", "email": "email", "created_at": "created_at"}
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
		return nil, 0, fmt.Errorf("list users: %w", err)
	}

	scanUser := func(rows *sql.Rows) (*model.User, error) {
		var u model.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
		return &u, err
	}
	users, err := database.ScanRows(rows, scanUser)
	if err != nil {
		return nil, 0, fmt.Errorf("scan users: %w", err)
	}
	return users, total, nil
}

func (r *UserRepository) ListByIDs(ctx context.Context, ids []int64) ([]*model.User, error) {
	if len(ids) == 0 {
		return []*model.User{}, nil
	}

	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE id IN (?" + strings.Repeat(",?", len(ids)-1) + ")"

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list users by ids: %w", err)
	}

	scanUser := func(rows *sql.Rows) (*model.User, error) {
		var u model.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
		return &u, err
	}

	users, err := database.ScanRows(rows, scanUser)
	if err != nil {
		return nil, fmt.Errorf("scan users by ids: %w", err)
	}

	return users, nil
}
