package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type APIGateway struct {
	router *mux.Router
	client *http.Client
}

type ServiceConfig struct {
	Name string
	URL  string
	Port string
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

func NewAPIGateway() *APIGateway {
	router := mux.NewRouter()

	// Configure CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Authorization", "Content-Type", "Accept"},
	})

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	gateway := &APIGateway{
		router: router,
		client: client,
	}

	// Setup routes
	gateway.setupRoutes()

	// Apply CORS middleware
	corsMiddleware.Handler(router)

	return &APIGateway{
		router: mux.NewRouter(),
		client: client,
	}
}

func (g *APIGateway) setupRoutes() {
	// Health check
	g.router.HandleFunc("/health", g.healthCheck).Methods("GET")

	// API v1 routes
	apiV1 := g.router.PathPrefix("/api/v1").Subrouter()

	// Product routes
	apiV1.HandleFunc("/products", g.handleProducts).Methods("GET")
	apiV1.HandleFunc("/products/{id}", g.handleProductByID).Methods("GET")

	// Payment routes
	apiV1.HandleFunc("/payments", g.handlePayments).Methods("POST")
	apiV1.HandleFunc("/payments/{id}", g.handlePaymentByID).Methods("GET")

	// Basket routes
	apiV1.HandleFunc("/baskets/{user_id}", g.handleBasket).Methods("GET")
	apiV1.HandleFunc("/baskets/add", g.handleAddItem).Methods("POST")
	apiV1.HandleFunc("/baskets/remove", g.handleRemoveItem).Methods("POST")

	// Metrics endpoint
	g.router.HandleFunc("/metrics", g.handleMetrics).Methods("GET")
}

func (g *APIGateway) healthCheck(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Success: true,
		Message: "API Gateway is healthy",
		Data: map[string]interface{}{
			"timestamp": time.Now().UTC(),
			"service":   "gingateway",
			"version":   "1.0.0",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (g *APIGateway) handleProducts(w http.ResponseWriter, r *http.Request) {
	// Forward to product service
	targetURL := "http://product-service:8081/v1/products"
	g.forwardRequest(w, r, targetURL)
}

func (g *APIGateway) handleProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	// Forward to product service
	targetURL := fmt.Sprintf("http://product-service:8081/v1/products/%s", productID)
	g.forwardRequest(w, r, targetURL)
}

func (g *APIGateway) handlePayments(w http.ResponseWriter, r *http.Request) {
	// Forward to payment service
	targetURL := "http://payment-service:8082/v1/payments"
	g.forwardRequest(w, r, targetURL)
}

func (g *APIGateway) handlePaymentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paymentID := vars["id"]

	// Forward to payment service
	targetURL := fmt.Sprintf("http://payment-service:8082/v1/payments/%s", paymentID)
	g.forwardRequest(w, r, targetURL)
}

func (g *APIGateway) handleBasket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// Forward to basket service
	targetURL := fmt.Sprintf("http://basket-service:8083/v1/baskets/%s", userID)
	g.forwardRequest(w, r, targetURL)
}

func (g *APIGateway) handleAddItem(w http.ResponseWriter, r *http.Request) {
	// Forward to basket service
	targetURL := "http://basket-service:8083/v1/baskets/add"
	g.forwardRequest(w, r, targetURL)
}

func (g *APIGateway) handleRemoveItem(w http.ResponseWriter, r *http.Request) {
	// Forward to basket service
	targetURL := "http://basket-service:8083/v1/baskets/remove"
	g.forwardRequest(w, r, targetURL)
}

func (g *APIGateway) handleMetrics(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Success: true,
		Data: map[string]interface{}{
			"uptime":    time.Since(startTime).String(),
			"requests":  requestCount,
			"timestamp": time.Now().UTC(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (g *APIGateway) forwardRequest(w http.ResponseWriter, r *http.Request, targetURL string) {
	// Increment request counter
	requestCount++

	// Parse target URL
	target, err := url.Parse(targetURL)
	if err != nil {
		g.sendError(w, "Invalid target URL", http.StatusBadRequest)
		return
	}

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Modify request
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
		req.URL.Path = target.Path
		req.URL.RawQuery = target.RawQuery

		// Add gateway headers
		req.Header.Set("X-Gateway", "gingateway")
		req.Header.Set("X-Forwarded-For", r.RemoteAddr)
		req.Header.Set("X-Request-ID", generateRequestID())
	}

	// Handle errors
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		g.sendError(w, "Service unavailable", http.StatusServiceUnavailable)
	}

	// Serve the request
	proxy.ServeHTTP(w, r)
}

func (g *APIGateway) sendError(w http.ResponseWriter, message string, statusCode int) {
	response := Response{
		Success: false,
		Error:   message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (g *APIGateway) Start(port string) error {
	log.Printf("Starting API Gateway on port %s", port)
	return http.ListenAndServe(":"+port, g.router)
}

// Global variables for metrics
var (
	startTime    = time.Now()
	requestCount = 0
)

func generateRequestID() string {
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}

func main() {
	port := getEnv("PORT", "8080")

	gateway := NewAPIGateway()

	log.Printf("API Gateway starting on port %s", port)
	if err := gateway.Start(port); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
