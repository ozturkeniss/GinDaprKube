package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"daprps/api/proto/events"
	productpb "daprps/api/proto/product"
	"daprps/internal/product-service/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductService struct {
	productpb.UnimplementedProductServiceServer
	repo model.ProductRepository
}

func NewProductService(repo model.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, req *productpb.GetProductRequest) (*productpb.GetProductResponse, error) {
	product, err := s.repo.GetByID(req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
	}

	return &productpb.GetProductResponse{
		Product: &productpb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Category:    product.Category,
		},
	}, nil
}

func (s *ProductService) UpdateStock(ctx context.Context, req *productpb.UpdateStockRequest) (*productpb.UpdateStockResponse, error) {
	updatedProduct, err := s.repo.UpdateStock(req.ProductId, req.Quantity, req.Operation)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating stock: %v", err)
	}

	return &productpb.UpdateStockResponse{
		Success:  true,
		NewStock: updatedProduct.Stock,
	}, nil
}

func (s *ProductService) ListProducts(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error) {
	products, err := s.repo.GetAll(req.Category, req.Limit, req.Offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error listing products: %v", err)
	}

	var protoProducts []*productpb.Product
	for _, p := range products {
		protoProducts = append(protoProducts, &productpb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			Category:    p.Category,
		})
	}

	return &productpb.ListProductsResponse{
		Products: protoProducts,
	}, nil
}

// Business logic methods
func (s *ProductService) CreateProduct(ctx context.Context, name, description, category string, price float64, stock int32) (*model.Product, error) {
	product := &model.Product{
		ID:          generateID(),
		Name:        name,
		Description: description,
		Category:    category,
		Price:       price,
		Stock:       stock,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repo.Create(product)
	if err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	return product, nil
}

func (s *ProductService) DecreaseStock(ctx context.Context, productID string, quantity int32) error {
	_, err := s.repo.UpdateStock(productID, quantity, "subtract")
	if err != nil {
		return fmt.Errorf("error decreasing stock: %w", err)
	}
	return nil
}

// Helper function to generate ID (in real app, use UUID)
func generateID() string {
	return fmt.Sprintf("prod_%d", time.Now().UnixNano())
}

// HandlePaymentCompleted implements PaymentEventHandler interface
func (s *ProductService) HandlePaymentCompleted(ctx context.Context, event *events.PaymentCompletedEvent) error {
	log.Printf("Received payment completed event for order %s, user %s", event.OrderId, event.UserId)

	// In a real application, you would:
	// 1. Get order details from the order service
	// 2. Update product stock based on ordered items
	// 3. Publish stock updated events

	// For now, we'll just log the event
	log.Printf("Payment completed for order %s with amount %f %s",
		event.OrderId, event.Amount, event.Currency)

	// TODO: Implement stock update logic when order service is available
	// This would involve:
	// - Getting order items from order service
	// - Updating product stock
	// - Publishing stock updated events

	return nil
}
