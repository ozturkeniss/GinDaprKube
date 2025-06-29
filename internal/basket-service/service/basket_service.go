package service

import (
	"context"
	"log"
	"time"

	basketpb "daprps/api/proto/basket"
	"daprps/api/proto/events"
	"daprps/internal/basket-service/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BasketService struct {
	basketpb.UnimplementedBasketServiceServer
	repo model.BasketRepository
}

func NewBasketService(repo model.BasketRepository) *BasketService {
	return &BasketService{
		repo: repo,
	}
}

func (s *BasketService) GetBasket(ctx context.Context, req *basketpb.GetBasketRequest) (*basketpb.GetBasketResponse, error) {
	basket, err := s.repo.GetByUserID(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting basket: %v", err)
	}

	return &basketpb.GetBasketResponse{
		Basket: &basketpb.Basket{
			UserId:      basket.UserID,
			Items:       convertBasketItems(basket.Items),
			TotalAmount: basket.TotalAmount,
			CreatedAt:   basket.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   basket.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *BasketService) AddItem(ctx context.Context, req *basketpb.AddItemRequest) (*basketpb.AddItemResponse, error) {
	// Note: In a real application, you would get product details from Product service
	// For now, we'll create a placeholder item
	item := model.BasketItem{
		ProductID:   req.ProductId,
		ProductName: "Product " + req.ProductId, // Placeholder
		Price:       0.0,                        // Will be updated from Product service
		Quantity:    req.Quantity,
	}

	err := s.repo.AddItem(req.UserId, item)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error adding item: %v", err)
	}

	// Get updated basket
	basket, err := s.repo.GetByUserID(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting updated basket: %v", err)
	}

	return &basketpb.AddItemResponse{
		Basket: &basketpb.Basket{
			UserId:      basket.UserID,
			Items:       convertBasketItems(basket.Items),
			TotalAmount: basket.TotalAmount,
			CreatedAt:   basket.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   basket.UpdatedAt.Format(time.RFC3339),
		},
		Success: true,
	}, nil
}

func (s *BasketService) RemoveItem(ctx context.Context, req *basketpb.RemoveItemRequest) (*basketpb.RemoveItemResponse, error) {
	err := s.repo.RemoveItem(req.UserId, req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error removing item: %v", err)
	}

	// Get updated basket
	basket, err := s.repo.GetByUserID(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting updated basket: %v", err)
	}

	return &basketpb.RemoveItemResponse{
		Basket: &basketpb.Basket{
			UserId:      basket.UserID,
			Items:       convertBasketItems(basket.Items),
			TotalAmount: basket.TotalAmount,
			CreatedAt:   basket.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   basket.UpdatedAt.Format(time.RFC3339),
		},
		Success: true,
	}, nil
}

func (s *BasketService) UpdateQuantity(ctx context.Context, req *basketpb.UpdateQuantityRequest) (*basketpb.UpdateQuantityResponse, error) {
	err := s.repo.UpdateQuantity(req.UserId, req.ProductId, req.Quantity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating quantity: %v", err)
	}

	// Get updated basket
	basket, err := s.repo.GetByUserID(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting updated basket: %v", err)
	}

	return &basketpb.UpdateQuantityResponse{
		Basket: &basketpb.Basket{
			UserId:      basket.UserID,
			Items:       convertBasketItems(basket.Items),
			TotalAmount: basket.TotalAmount,
			CreatedAt:   basket.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   basket.UpdatedAt.Format(time.RFC3339),
		},
		Success: true,
	}, nil
}

func (s *BasketService) ClearBasket(ctx context.Context, req *basketpb.ClearBasketRequest) (*basketpb.ClearBasketResponse, error) {
	err := s.repo.Clear(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error clearing basket: %v", err)
	}

	return &basketpb.ClearBasketResponse{
		Success: true,
	}, nil
}

// Business logic methods
func (s *BasketService) GetBasketByUserID(ctx context.Context, userID string) (*model.Basket, error) {
	return s.repo.GetByUserID(userID)
}

func (s *BasketService) ClearBasketByUserID(ctx context.Context, userID string) error {
	return s.repo.Clear(userID)
}

// HandlePaymentCompleted implements PaymentEventHandler interface
func (s *BasketService) HandlePaymentCompleted(ctx context.Context, event *events.PaymentCompletedEvent) error {
	log.Printf("Received payment completed event for order %s, user %s", event.OrderId, event.UserId)

	// Clear the user's basket after successful payment
	err := s.repo.Clear(event.UserId)
	if err != nil {
		log.Printf("Failed to clear basket for user %s: %v", event.UserId, err)
		return err
	}

	log.Printf("Successfully cleared basket for user %s after payment completion", event.UserId)
	return nil
}

// Helper functions
func convertBasketItems(items []model.BasketItem) []*basketpb.BasketItem {
	var protoItems []*basketpb.BasketItem
	for _, item := range items {
		protoItems = append(protoItems, &basketpb.BasketItem{
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Price:       item.Price,
			Quantity:    item.Quantity,
		})
	}
	return protoItems
}
