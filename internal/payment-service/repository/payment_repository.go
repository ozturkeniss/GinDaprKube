package repository

import (
	"fmt"
	"time"

	"daprps/internal/payment-service/model"

	"gorm.io/gorm"
)

type PaymentRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) model.PaymentRepository {
	return &PaymentRepositoryImpl{db: db}
}

func (r *PaymentRepositoryImpl) GetByID(id string) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Where("id = ?", id).First(&payment).Error
	if err != nil {
		return nil, fmt.Errorf("error getting payment by ID: %w", err)
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) GetByOrderID(orderID string) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		return nil, fmt.Errorf("error getting payment by order ID: %w", err)
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) Create(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepositoryImpl) Update(payment *model.Payment) error {
	payment.UpdatedAt = time.Now()
	return r.db.Save(payment).Error
}

func (r *PaymentRepositoryImpl) UpdateStatus(id, status string) error {
	return r.db.Model(&model.Payment{}).Where("id = ?", id).Update("status", status).Error
}
