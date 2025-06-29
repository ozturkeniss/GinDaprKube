package model

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID            string         `json:"id" gorm:"primaryKey;type:varchar(255)"`
	OrderID       string         `json:"order_id" gorm:"type:varchar(255);not null"`
	UserID        string         `json:"user_id" gorm:"type:varchar(255)"`
	Amount        float64        `json:"amount"`
	Currency      string         `json:"currency" gorm:"type:varchar(10)"`
	Status        string         `json:"status" gorm:"type:varchar(50);not null"`
	PaymentMethod string         `json:"payment_method" gorm:"type:varchar(50)"`
	CardNumber    string         `json:"card_number" gorm:"type:varchar(32)"`
	CardHolder    string         `json:"card_holder" gorm:"type:varchar(100)"`
	ExpiryDate    string         `json:"expiry_date" gorm:"type:varchar(10)"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type PaymentRepository interface {
	GetByID(id string) (*Payment, error)
	GetByOrderID(orderID string) (*Payment, error)
	Create(payment *Payment) error
	Update(payment *Payment) error
	UpdateStatus(id, status string) error
}
