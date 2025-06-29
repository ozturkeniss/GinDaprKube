package model

import (
	"time"
)

type BasketItem struct {
	ProductID   string  `json:"product_id" db:"product_id"`
	ProductName string  `json:"product_name" db:"product_name"`
	Price       float64 `json:"price" db:"price"`
	Quantity    int32   `json:"quantity" db:"quantity"`
}

type Basket struct {
	UserID      string       `json:"user_id" db:"user_id"`
	Items       []BasketItem `json:"items" db:"items"`
	TotalAmount float64      `json:"total_amount" db:"total_amount"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}

type BasketRepository interface {
	GetByUserID(userID string) (*Basket, error)
	AddItem(userID string, item BasketItem) error
	RemoveItem(userID, productID string) error
	UpdateQuantity(userID, productID string, quantity int32) error
	Clear(userID string) error
	UpdateTotalAmount(userID string, totalAmount float64) error
}
