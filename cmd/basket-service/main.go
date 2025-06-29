package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"

	"daprps/api/proto/basket"
	"daprps/internal/basket-service/repository"
	"daprps/internal/basket-service/service"
	"daprps/kafka/consumer"
)

func main() {
	// Get Redis configuration from environment variables
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDBStr := getEnv("REDIS_DB", "0")

	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		log.Printf("Invalid REDIS_DB value, using 0: %v", err)
		redisDB = 0
	}

	// Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Test Redis connection
	ctx := rdb.Context()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Redis not available, using in-memory storage: %v", err)
		// In a real app, you might want to exit here
	}

	// Create repository and service
	repo := repository.NewBasketRepository(redisHost+":"+redisPort, redisPassword, redisDB)
	basketService := service.NewBasketService(repo)

	// Create Kafka consumer
	kafkaConsumer, err := consumer.NewPaymentConsumer(basketService)
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
	basket.RegisterBasketServiceServer(grpcServer, basketService)

	// Start gRPC server
	grpcPort := getEnv("GRPC_PORT", "50053")
	grpcLis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port: %v", err)
	}

	go func() {
		log.Printf("Basket service gRPC starting on :%s", grpcPort)
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Create HTTP server for API Gateway
	httpPort := getEnv("HTTP_PORT", "8083")
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// Get basket endpoint
	mux.HandleFunc("/v1/baskets/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract user ID from URL
		userID := r.URL.Path[len("/v1/baskets/"):]
		if userID == "" {
			http.Error(w, "User ID required", http.StatusBadRequest)
			return
		}

		// Get basket
		basket, err := basketService.GetBasket(context.Background(), &basket.GetBasketRequest{UserId: userID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(basket)
	})

	// Add item to basket endpoint
	mux.HandleFunc("/v1/baskets/add", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			UserID    string `json:"user_id"`
			ProductID string `json:"product_id"`
			Quantity  int32  `json:"quantity"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Add item to basket
		basket, err := basketService.AddItem(context.Background(), &basket.AddItemRequest{
			UserId:    req.UserID,
			ProductId: req.ProductID,
			Quantity:  req.Quantity,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(basket)
	})

	// Remove item from basket endpoint
	mux.HandleFunc("/v1/baskets/remove", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			UserID    string `json:"user_id"`
			ProductID string `json:"product_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Remove item from basket
		basket, err := basketService.RemoveItem(context.Background(), &basket.RemoveItemRequest{
			UserId:    req.UserID,
			ProductId: req.ProductID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(basket)
	})

	log.Printf("Basket service HTTP starting on :%s", httpPort)
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
