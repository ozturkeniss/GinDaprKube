package repository

import (
	"fmt"
	"time"

	"daprps/internal/product-service/model"

	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) model.ProductRepository {
	return &ProductRepositoryImpl{db: db}
}

func (r *ProductRepositoryImpl) GetByID(id string) (*model.Product, error) {
	var product model.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, fmt.Errorf("error getting product by ID: %w", err)
	}
	return &product, nil
}

func (r *ProductRepositoryImpl) GetAll(category string, limit, offset int32) ([]*model.Product, error) {
	var products []*model.Product
	query := r.db.Order("created_at DESC")

	if category != "" {
		query = query.Where("category = ?", category)
	}

	err := query.Limit(int(limit)).Offset(int(offset)).Find(&products).Error
	if err != nil {
		return nil, fmt.Errorf("error getting products: %w", err)
	}

	return products, nil
}

func (r *ProductRepositoryImpl) UpdateStock(id string, quantity int32, operation string) (*model.Product, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("error starting transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get current product
	var product model.Product
	err := tx.Where("id = ?", id).First(&product).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error getting product: %w", err)
	}

	// Calculate new stock
	switch operation {
	case "add":
		product.Stock += quantity
	case "subtract":
		product.Stock -= quantity
		if product.Stock < 0 {
			tx.Rollback()
			return nil, fmt.Errorf("insufficient stock")
		}
	default:
		tx.Rollback()
		return nil, fmt.Errorf("invalid operation: %s", operation)
	}

	product.UpdatedAt = time.Now()

	// Update product
	err = tx.Save(&product).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error updating product: %w", err)
	}

	if err = tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &product, nil
}

func (r *ProductRepositoryImpl) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepositoryImpl) Update(product *model.Product) error {
	product.UpdatedAt = time.Now()
	return r.db.Save(product).Error
}

func (r *ProductRepositoryImpl) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Product{}).Error
}
