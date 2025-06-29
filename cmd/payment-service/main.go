package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"daprps/api/proto/payment"
	"daprps/internal/payment-service/model"
	"daprps/internal/payment-service/repository"
	"daprps/internal/payment-service/service"
	"daprps/kafka/publisher"
)

func main() {
	// Get database configuration from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "paymentdb")

	// Database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migration
	err = db.AutoMigrate(&model.Payment{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated successfully")

	// Create Kafka publisher
	kafkaPublisher, err := publisher.NewPaymentPublisher()
	if err != nil {
		log.Fatalf("Failed to create Kafka publisher: %v", err)
	}
	defer kafkaPublisher.Close()

	// Create repository and service
	repo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(repo, kafkaPublisher)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	payment.RegisterPaymentServiceServer(grpcServer, paymentService)

	// Start gRPC server
	grpcPort := getEnv("GRPC_PORT", "50052")
	grpcLis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port: %v", err)
	}

	go func() {
		log.Printf("Payment service gRPC starting on :%s", grpcPort)
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Create HTTP server for GinGateway
	httpPort := getEnv("HTTP_PORT", "8082")
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// Process payment endpoint
	mux.HandleFunc("/v1/payments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse request body
		var req payment.ProcessPaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Process payment
		resp, err := paymentService.ProcessPayment(context.Background(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(resp)
	})

	// Get payment status endpoint
	mux.HandleFunc("/v1/payments/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract payment ID from URL
		paymentID := r.URL.Path[len("/v1/payments/"):]
		if paymentID == "" {
			http.Error(w, "Payment ID required", http.StatusBadRequest)
			return
		}

		// Get payment status
		resp, err := paymentService.GetPaymentStatus(context.Background(), &payment.GetPaymentStatusRequest{PaymentId: paymentID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(resp)
	})

	log.Printf("Payment service HTTP starting on :%s", httpPort)
	if err := http.ListenAndServe(":"+httpPort, mux); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
