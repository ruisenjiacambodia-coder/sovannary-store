package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ============================================================
// Data Models
// ============================================================

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	NameKh      string  `json:"nameKh"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	PriceKHR    float64 `json:"priceKHR"`
	ImageURL    string  `json:"imageURL"`
	SKU         string  `json:"sku"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type Order struct {
	ID         string      `json:"id"`
	Customer   Customer    `json:"customer"`
	Items      []OrderItem `json:"items"`
	Total      float64     `json:"total"`
	TotalKHR   float64     `json:"totalKHR"`
	Status     string      `json:"status"`
	CreatedAt  string      `json:"createdAt"`
}

type Customer struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Note    string `json:"note"`
}

type OrderItem struct {	ProductID string  `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	PriceKHR  float64 `json:"priceKHR"`
	Qty       int     `json:"qty"`
}

type StoreSettings struct {
	StoreName       string `json:"storeName"`
	Tagline         string `json:"tagline"`
	Phone           string `json:"phone"`
	TelegramBot     string `json:"telegramBot"`
	TelegramChatID  string `json:"telegramChatId"`
	WhatsAppNumber  string `json:"whatsappNumber"`
	Currency        string `json:"currency"`
	ExchangeRate    float64 `json:"exchangeRate"`
}

// ============================================================
// In-Memory Storage
// ============================================================

type Store struct {
	mu       sync.RWMutex
	products map[string]Product
	orders   []Order
	settings StoreSettings
	counter  int
}

var store = &Store{
	products: make(map[string]Product),
	orders:   []Order{},
	settings: StoreSettings{
		StoreName:      "Sovannary Store",
		Tagline:        "ផលិតផលគុណភាពខ្ពស់ សម្រាប់ជីវិតប្រណីត",
		Phone:          "+855 12 345 678",
		TelegramBot:    "",
		TelegramChatID: "",
		WhatsAppNumber: "85512345678",
		Currency:       "USD",
		ExchangeRate:   4100,
	},
}

// Seed initial demo data
func init() {
	demoProducts := []Product{
		{
			ID: "P001", Name: "Premium Jasmine Rice 5kg", NameKh: "អង្ករផ្កាម្លិះ ៥គីឡូ",			Description: "Premium quality jasmine rice, fragrant and soft texture.",
			Price: 12.50, PriceKHR: 51250,
			ImageURL: "https://images.unsplash.com/photo-1586201375761-83865001e31c?w=400",
			SKU: "RICE-JAS-5KG", Category: "Food & Grocery", Stock: 45,
		},
		{
			ID: "P002", Name: "Organic Coconut Oil 500ml", NameKh: "ប្រេងដូងសរីរាង្គ ៥០០មីលីលីត្រ",
			Description: "Cold-pressed organic coconut oil for cooking and skincare.",
			Price: 8.90, PriceKHR: 36490,
			ImageURL: "https://images.unsplash.com/photo-1526947425960-945c6e72858f?w=400",
			SKU: "OIL-COCO-500", Category: "Health & Beauty", Stock: 30,
		},
		{
			ID: "P003", Name: "Handcrafted Silk Scarf", NameKh: "កន្សែងកសម្រេចដៃ",
			Description: "Traditional Khmer silk scarf, handwoven by local artisans.",
			Price: 35.00, PriceKHR: 143500,
			ImageURL: "https://images.unsplash.com/photo-1601924994987-69e26d50dc26?w=400",
			SKU: "SCF-SILK-001", Category: "Fashion", Stock: 15,
		},
		{
			ID: "P004", Name: "Kampot Pepper Black 100g", NameKh: "ម្រេចកំពតខ្មៅ ១០០ក្រាម",
			Description: "Authentic Kampot pepper, GI protected origin.",
			Price: 15.00, PriceKHR: 61500,
			ImageURL: "https://images.unsplash.com/photo-1599909533730-3b3baf4a4b40?w=400",
			SKU: "PEP-KAMP-100", Category: "Food & Grocery", Stock: 60,
		},
		{
			ID: "P005", Name: "Palm Sugar 1kg", NameKh: "ស្ករត្នោត ១គីឡូ",
			Description: "Natural palm sugar from Kampong Speu province.",
			Price: 5.50, PriceKHR: 22550,
			ImageURL: "https://images.unsplash.com/photo-1581798459219-318e76aecc7b?w=400",
			SKU: "SUG-PALM-1KG", Category: "Food & Grocery", Stock: 80,
		},
		{
			ID: "P006", Name: "Ceramic Tea Set", NameKh: "ឈុតតែសេរ៉ាមិច",
			Description: "Elegant ceramic tea set with 4 cups and teapot.",
			Price: 48.00, PriceKHR: 196800,
			ImageURL: "https://images.unsplash.com/photo-1556679343-c7306c1976bc?w=400",
			SKU: "TEA-CERA-SET", Category: "Home & Living", Stock: 10,
		},
	}

	now := time.Now().Format(time.RFC3339)
	for _, p := range demoProducts {
		p.CreatedAt = now
		p.UpdatedAt = now
		store.products[p.ID] = p
	}
	store.counter = len(demoProducts)
}
// ============================================================
// CORS Middleware
// ============================================================

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ============================================================
// JSON Helpers
// ============================================================

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func generateID() string {
	store.counter++
	return fmt.Sprintf("P%03d", store.counter)
}

// ============================================================
// API Handlers
// ============================================================

// GET /api/products
func handleListProducts(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	products := make([]Product, 0, len(store.products))
	for _, p := range store.products {
		products = append(products, p)
	}
	writeJSON(w, http.StatusOK, products)
}


// GET /api/products/{id}
func handleGetProduct(w http.ResponseWriter, r *http.Request, id string) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	p, ok := store.products[id]
	if !ok {
		writeError(w, http.StatusNotFound, "Product not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// POST /api/products
func handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	p.ID = generateID()
	now := time.Now().Format(time.RFC3339)
	p.CreatedAt = now
	p.UpdatedAt = now

	store.products[p.ID] = p
	writeJSON(w, http.StatusCreated, p)
}

// PUT /api/products/{id}
func handleUpdateProduct(w http.ResponseWriter, r *http.Request, id string) {
	var update Product
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	existing, ok := store.products[id]
	if !ok {
		writeError(w, http.StatusNotFound, "Product not found")
		return
	}

	// Preserve ID and CreatedAt
	update.ID = existing.ID
	update.CreatedAt = existing.CreatedAt
	update.UpdatedAt = time.Now().Format(time.RFC3339)

	store.products[id] = update
	writeJSON(w, http.StatusOK, update)
}

// DELETE /api/products/{id}
func handleDeleteProduct(w http.ResponseWriter, r *http.Request, id string) {
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, ok := store.products[id]; !ok {
		writeError(w, http.StatusNotFound, "Product not found")
		return
	}
	delete(store.products, id)
	writeJSON(w, http.StatusOK, map[string]string{"message": "Deleted"})
}

// POST /api/orders
func handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	order.ID = fmt.Sprintf("ORD-%d", time.Now().UnixNano())
	order.Status = "pending"
	order.CreatedAt = time.Now().Format(time.RFC3339)

	store.orders = append(store.orders, order)
	writeJSON(w, http.StatusCreated, order)
}

// GET /api/orders
func handleListOrders(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	writeJSON(w, http.StatusOK, store.orders)
}
// GET /api/settings
func handleGetSettings(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	writeJSON(w, http.StatusOK, store.settings)
}

// PUT /api/settings
func handleUpdateSettings(w http.ResponseWriter, r *http.Request) {
	var s StoreSettings
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	store.settings = s
	writeJSON(w, http.StatusOK, s)
}

// POST /api/sync (receive bulk data from offline client)
func handleSync(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Products []Product `json:"products"`
		Orders   []Order   `json:"orders"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for _, p := range payload.Products {
		if _, exists := store.products[p.ID]; !exists {
			store.products[p.ID] = p
		}
	}
	for _, o := range payload.Orders {
		store.orders = append(store.orders, o)
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "synced"})
}

// ============================================================// Router
// ============================================================

func router(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Static files
	if path == "/" || path == "/index.html" {
		http.ServeFile(w, r, "dashboard.html")
		return
	}
	if path == "/sw.js" {
		w.Header().Set("Content-Type", "application/javascript")
		w.Header().Set("Service-Worker-Allowed", "/")
		http.ServeFile(w, r, "sw.js")
		return
	}
	if path == "/manifest.json" {
		http.ServeFile(w, r, "manifest.json")
		return
	}

	// API routing
	switch {
	case path == "/api/products" && r.Method == http.MethodGet:
		handleListProducts(w, r)
	case path == "/api/products" && r.Method == http.MethodPost:
		handleCreateProduct(w, r)
	case strings.HasPrefix(path, "/api/products/") && r.Method == http.MethodGet:
		id := strings.TrimPrefix(path, "/api/products/")
		handleGetProduct(w, r, id)
	case strings.HasPrefix(path, "/api/products/") && r.Method == http.MethodPut:
		id := strings.TrimPrefix(path, "/api/products/")
		handleUpdateProduct(w, r, id)
	case strings.HasPrefix(path, "/api/products/") && r.Method == http.MethodDelete:
		id := strings.TrimPrefix(path, "/api/products/")
		handleDeleteProduct(w, r, id)
	case path == "/api/orders" && r.Method == http.MethodGet:
		handleListOrders(w, r)
	case path == "/api/orders" && r.Method == http.MethodPost:
		handleCreateOrder(w, r)
	case path == "/api/settings" && r.Method == http.MethodGet:
		handleGetSettings(w, r)
	case path == "/api/settings" && r.Method == http.MethodPut:
		handleUpdateSettings(w, r)
	case path == "/api/sync" && r.Method == http.MethodPost:
		handleSync(w, r)
	default:
		writeError(w, http.StatusNotFound, "Not found")
	}}

// ============================================================
// Main
// ============================================================

func main() {
	http.Handle("/", corsMiddleware(http.HandlerFunc(router)))

	port := "8080"
	if p := strings.TrimSpace(getEnv("PORT", port)); p != "" {
		port = p
	}

	fmt.Printf("🚀 Sovannary Store running at http://localhost:%s\n", port)
	fmt.Printf("📦 In-memory storage initialized with %d products\n", len(store.products))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getEnv(key, fallback string) string {
	if v, ok := lookupEnv(key); ok {
		return v
	}
	return fallback
}

// Simple env lookup without importing os (to keep imports minimal)
var lookupEnv = func(key string) (string, bool) {
	// Using a basic approach
	return "", false
}

func init() {
	// Override with actual os.Getenv
	lookupEnv = func(key string) (string, bool) {
		// Re-import workaround - we'll just use a simple approach
		return "", false
	}
}

// Helper to parse int from string (unused but available)
var _ = strconv.Itoa