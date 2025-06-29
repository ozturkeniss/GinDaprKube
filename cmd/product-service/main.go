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

	"daprps/api/proto/product"
	"daprps/internal/product-service/model"
	"daprps/internal/product-service/repository"
	"daprps/internal/product-service/service"
	"daprps/kafka/consumer"
)

func main() {
	// Get database configuration from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "productdb")

	// Database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migration
	err = db.AutoMigrate(&model.Product{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated successfully")

	// Create repository and service
	repo := repository.NewProductRepository(db)
	productService := service.NewProductService(repo)

	// Create Kafka consumer for payment events
	kafkaConsumer, err := consumer.NewPaymentConsumer(productService)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer kafkaConsumer.Close()

	// Start Kafka consumer in background
	go func() {
		topics := []string{"payment-completed"}
		if err := kafkaConsumer.Start(context.Background(), topics); err != nil {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	// Create gRPC server
	grpcServer := grpc.NewServer()
	product.RegisterProductServiceServer(grpcServer, productService)

	// Start gRPC server
	grpcPort := getEnv("GRPC_PORT", "50051")
	grpcLis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port: %v", err)
	}

	go func() {
		log.Printf("Product service gRPC starting on :%s", grpcPort)
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Create HTTP server for KrakenD
	httpPort := getEnv("HTTP_PORT", "8081")
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// Products endpoint
	mux.HandleFunc("/v1/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get all products
		products, err := productService.ListProducts(context.Background(), &product.ListProductsRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(products)
	})

	// Product by ID endpoint
	mux.HandleFunc("/v1/products/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract product ID from URL
		productID := r.URL.Path[len("/v1/products/"):]
		if productID == "" {
			http.Error(w, "Product ID required", http.StatusBadRequest)
			return
		}

		// Get product by ID
		product, err := productService.GetProduct(context.Background(), &product.GetProductRequest{ProductId: productID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(product)
	})

	log.Printf("Product service HTTP starting on :%s", httpPort)
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
