package service

import (
	"context"
	"fmt"
	"time"

	"daprps/api/proto/events"
	paymentpb "daprps/api/proto/payment"
	"daprps/internal/payment-service/model"
	"daprps/kafka/publisher"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentService struct {
	paymentpb.UnimplementedPaymentServiceServer
	repo      model.PaymentRepository
	publisher *publisher.PaymentPublisher
}

func NewPaymentService(repo model.PaymentRepository, publisher *publisher.PaymentPublisher) *PaymentService {
	return &PaymentService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, req *paymentpb.ProcessPaymentRequest) (*paymentpb.ProcessPaymentResponse, error) {
	// Create payment record
	payment := &model.Payment{
		ID:            generatePaymentID(),
		OrderID:       req.OrderId,
		UserID:        "", // Will be set from context or request
		Amount:        req.Amount,
		Currency:      req.Currency,
		Status:        "pending",
		PaymentMethod: req.PaymentMethod,
		CardNumber:    req.CardNumber,
		CardHolder:    req.CardHolder,
		ExpiryDate:    req.ExpiryDate,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save payment to database
	err := s.repo.Create(payment)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating payment: %v", err)
	}

	// Simulate payment processing
	// In real app, this would call payment gateway
	time.Sleep(100 * time.Millisecond)

	// Update payment status to completed
	err = s.repo.UpdateStatus(payment.ID, "completed")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating payment status: %v", err)
	}

	// Publish payment completed event
	event := &events.PaymentCompletedEvent{
		PaymentId:     payment.ID,
		OrderId:       payment.OrderID,
		UserId:        payment.UserID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		PaymentMethod: payment.PaymentMethod,
		CompletedAt:   time.Now().Format(time.RFC3339),
	}

	err = s.publisher.PublishPaymentCompleted(ctx, event)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to publish payment completed event: %v\n", err)
	}

	// Get updated payment
	updatedPayment, err := s.repo.GetByID(payment.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting updated payment: %v", err)
	}

	return &paymentpb.ProcessPaymentResponse{
		Payment: &paymentpb.Payment{
			Id:            updatedPayment.ID,
			OrderId:       updatedPayment.OrderID,
			Amount:        updatedPayment.Amount,
			Currency:      updatedPayment.Currency,
			Status:        updatedPayment.Status,
			PaymentMethod: updatedPayment.PaymentMethod,
			CreatedAt:     updatedPayment.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     updatedPayment.UpdatedAt.Format(time.RFC3339),
		},
		Success: true,
	}, nil
}

func (s *PaymentService) GetPaymentStatus(ctx context.Context, req *paymentpb.GetPaymentStatusRequest) (*paymentpb.GetPaymentStatusResponse, error) {
	payment, err := s.repo.GetByID(req.PaymentId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
	}

	return &paymentpb.GetPaymentStatusResponse{
		Payment: &paymentpb.Payment{
			Id:            payment.ID,
			OrderId:       payment.OrderID,
			Amount:        payment.Amount,
			Currency:      payment.Currency,
			Status:        payment.Status,
			PaymentMethod: payment.PaymentMethod,
			CreatedAt:     payment.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     payment.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *PaymentService) RefundPayment(ctx context.Context, req *paymentpb.RefundPaymentRequest) (*paymentpb.RefundPaymentResponse, error) {
	// Get original payment
	payment, err := s.repo.GetByID(req.PaymentId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
	}

	// Check if payment is completed
	if payment.Status != "completed" {
		return nil, status.Errorf(codes.FailedPrecondition, "payment is not completed, cannot refund")
	}

	// Simulate refund processing
	// In real app, this would call payment gateway
	time.Sleep(100 * time.Millisecond)

	// Update payment status to refunded
	err = s.repo.UpdateStatus(payment.ID, "refunded")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating payment status: %v", err)
	}

	return &paymentpb.RefundPaymentResponse{
		Success:  true,
		RefundId: generateRefundID(),
	}, nil
}

// Business logic methods
func (s *PaymentService) ProcessPaymentWithItems(ctx context.Context, orderID, userID string, amount float64, currency, paymentMethod string) (*model.Payment, error) {
	payment := &model.Payment{
		ID:            generatePaymentID(),
		OrderID:       orderID,
		UserID:        userID,
		Amount:        amount,
		Currency:      currency,
		Status:        "pending",
		PaymentMethod: paymentMethod,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := s.repo.Create(payment)
	if err != nil {
		return nil, fmt.Errorf("error creating payment: %w", err)
	}

	// Simulate payment processing
	time.Sleep(100 * time.Millisecond)

	// Update to completed
	err = s.repo.UpdateStatus(payment.ID, "completed")
	if err != nil {
		return nil, fmt.Errorf("error updating payment status: %w", err)
	}

	return payment, nil
}

// Helper functions
func generatePaymentID() string {
	return fmt.Sprintf("pay_%d", time.Now().UnixNano())
}

func generateRefundID() string {
	return fmt.Sprintf("ref_%d", time.Now().UnixNano())
}
