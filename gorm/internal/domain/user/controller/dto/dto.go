package dto

import "db_blueprints/pkgs/paging"

type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ListUserRequest struct {
	Search    string `json:"search,omitempty" form:"search"`
	Page      int64  `json:"-" form:"page"`
	Limit     int64  `json:"-" form:"size"`
	OrderBy   string `json:"-" form:"order_by"`
	OrderDesc bool   `json:"-" form:"order_desc"`
	TakeAll   bool   `json:"-" form:"take_all"`
}

type ListUserResponse struct {
	Users      []*User            `json:"items"`
	Pagination *paging.Pagination `json:"metadata"`
}

type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type CreateUserResponse struct {
	User *User `json:"user"`
}

type UpdateUserRequest struct {
	ID    int64   `json:"id"`
	Email *string `json:"email"`
	Name  *string `json:"name"`
}

type UpdateUserResponse struct {
	User *User `json:"user"`
}
