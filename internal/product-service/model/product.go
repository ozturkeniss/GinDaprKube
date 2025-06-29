package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          string         `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name        string         `json:"name" gorm:"type:varchar(255);not null"`
	Description string         `json:"description" gorm:"type:text"`
	Price       float64        `json:"price"`
	Stock       int32          `json:"stock" gorm:"type:int;not null;default:0"`
	Category    string         `json:"category" gorm:"type:varchar(100)"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type ProductRepository interface {
	GetByID(id string) (*Product, error)
	GetAll(category string, limit, offset int32) ([]*Product, error)
	UpdateStock(id string, quantity int32, operation string) (*Product, error)
	Create(product *Product) error
	Update(product *Product) error
	Delete(id string) error
}
