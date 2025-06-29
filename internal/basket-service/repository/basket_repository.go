package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"daprps/internal/basket-service/model"

	"github.com/go-redis/redis/v8"
)

type BasketRepositoryImpl struct {
	client *redis.Client
}

func NewBasketRepository(addr, password string, db int) model.BasketRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &BasketRepositoryImpl{
		client: client,
	}
}

func (r *BasketRepositoryImpl) GetByUserID(userID string) (*model.Basket, error) {
	ctx := context.Background()
	key := fmt.Sprintf("basket:%s", userID)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// Basket doesn't exist, return empty basket
			return &model.Basket{
				UserID:      userID,
				Items:       []model.BasketItem{},
				TotalAmount: 0,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}, nil
		}
		return nil, fmt.Errorf("error getting basket: %w", err)
	}

	var basket model.Basket
	err = json.Unmarshal([]byte(data), &basket)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling basket: %w", err)
	}

	return &basket, nil
}

func (r *BasketRepositoryImpl) AddItem(userID string, item model.BasketItem) error {
	basket, err := r.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("error getting basket: %w", err)
	}

	// Check if item already exists
	itemExists := false
	for i, existingItem := range basket.Items {
		if existingItem.ProductID == item.ProductID {
			basket.Items[i].Quantity += item.Quantity
			itemExists = true
			break
		}
	}

	if !itemExists {
		basket.Items = append(basket.Items, item)
	}

	// Recalculate total
	basket.TotalAmount = 0
	for _, item := range basket.Items {
		basket.TotalAmount += item.Price * float64(item.Quantity)
	}

	basket.UpdatedAt = time.Now()

	return r.saveBasket(basket)
}

func (r *BasketRepositoryImpl) RemoveItem(userID, productID string) error {
	basket, err := r.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("error getting basket: %w", err)
	}

	// Remove item
	var newItems []model.BasketItem
	for _, item := range basket.Items {
		if item.ProductID != productID {
			newItems = append(newItems, item)
		}
	}

	basket.Items = newItems

	// Recalculate total
	basket.TotalAmount = 0
	for _, item := range basket.Items {
		basket.TotalAmount += item.Price * float64(item.Quantity)
	}

	basket.UpdatedAt = time.Now()

	return r.saveBasket(basket)
}

func (r *BasketRepositoryImpl) UpdateQuantity(userID, productID string, quantity int32) error {
	basket, err := r.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("error getting basket: %w", err)
	}

	// Update quantity
	for i, item := range basket.Items {
		if item.ProductID == productID {
			if quantity <= 0 {
				// Remove item if quantity is 0 or negative
				basket.Items = append(basket.Items[:i], basket.Items[i+1:]...)
			} else {
				basket.Items[i].Quantity = quantity
			}
			break
		}
	}

	// Recalculate total
	basket.TotalAmount = 0
	for _, item := range basket.Items {
		basket.TotalAmount += item.Price * float64(item.Quantity)
	}

	basket.UpdatedAt = time.Now()

	return r.saveBasket(basket)
}

func (r *BasketRepositoryImpl) Clear(userID string) error {
	basket, err := r.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("error getting basket: %w", err)
	}

	basket.Items = []model.BasketItem{}
	basket.TotalAmount = 0
	basket.UpdatedAt = time.Now()

	return r.saveBasket(basket)
}

func (r *BasketRepositoryImpl) UpdateTotalAmount(userID string, totalAmount float64) error {
	basket, err := r.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("error getting basket: %w", err)
	}

	basket.TotalAmount = totalAmount
	basket.UpdatedAt = time.Now()

	return r.saveBasket(basket)
}

func (r *BasketRepositoryImpl) saveBasket(basket *model.Basket) error {
	ctx := context.Background()
	key := fmt.Sprintf("basket:%s", basket.UserID)

	data, err := json.Marshal(basket)
	if err != nil {
		return fmt.Errorf("error marshaling basket: %w", err)
	}

	err = r.client.Set(ctx, key, data, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("error saving basket: %w", err)
	}

	return nil
}

func (r *BasketRepositoryImpl) Close() error {
	return r.client.Close()
}
