package model

import "time"

// User defines the model for a user in the system.
// The `json`, `db`, and `gorm` struct tags are used for compatibility
// with JSON encoding, the sqlx library, and the GORM ORM, respectively.
type User struct {
	ID        int64     `json:"id" db:"id" gorm:"column:id;primaryKey"`
	Name      string    `json:"name" db:"name" gorm:"column:name"`
	Email     string    `json:"email" db:"email" gorm:"column:email"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"column:updated_at"`

	// Products represents the one-to-many relationship: a User has many Products.
	// This field is primarily used by GORM for preloading/joining.
	Products []Product `json:"-" gorm:"foreignKey:OwnerID"`
}

// TableName explicitly specifies the table name for GORM.
func (u *User) TableName() string {
	u.CreatedAt = time.Now().UTC()
	u.UpdatedAt = time.Now().UTC()

	return "users"
}
