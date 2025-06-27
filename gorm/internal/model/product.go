package model

import "time"

// Product defines the product model.
// It has a "belongs to" relationship with a User.
type Product struct {
	ID        int64     `json:"id" db:"id" gorm:"column:id;primaryKey"`
	Name      string    `json:"name" db:"name" gorm:"column:name"`
	Price     float64   `json:"price" db:"price" gorm:"column:price"` // Note: For real-world financial applications, consider using a decimal type or storing as an integer (in cents).
	OwnerID   int64     `json:"owner_id" db:"owner_id" gorm:"column:owner_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"column:updated_at"`

	// Owner represents the many-to-one relationship: a Product belongs to one User.
	// This field is used by GORM to preload/join the owner's data.
	Owner User `json:"user" gorm:"foreignKey:OwnerID"`
}

// TableName explicitly specifies the table name for GORM.
func (p *Product) TableName() string {
	p.CreatedAt = time.Now().UTC()
	p.UpdatedAt = time.Now().UTC()

	return "products"
}
