package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"clientsweb-backend/internal/handler"
	"clientsweb-backend/internal/service"
	"clientsweb-backend/internal/store"
)

func main() {
	// Load configuration from config.json
	cfg, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize JSON store with data file path
	customerStore, err := store.NewJSONStore("data/customers.json")
	if err != nil {
		log.Fatalf("Failed to initialize customer store: %v", err)
	}

	// Initialize license store
	licenseStore, err := store.NewLicenseStore("data/licenses.json")
	if err != nil {
		log.Fatalf("Failed to initialize license store: %v", err)
	}

	// Create service and handler layers
	customerService := service.NewCustomerService(customerStore, licenseStore)
	customerHandler := handler.NewCustomerHandler(customerService)

	// Determine frontend directory path (works whether run from backend/ or elsewhere)
	frontendDir := "../frontend"
	if _, err := os.Stat(frontendDir + "/index.html"); os.IsNotExist(err) {
		// Fallback: try ./frontend if ../frontend doesn't exist
		if _, err2 := os.Stat("./frontend/index.html"); err2 == nil {
			frontendDir = "./frontend"
		}
	}

	// Set up router with chi
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware) // Allow frontend to access API

	// Register routes under /api/customers
	customerRouter := customerHandler.Routes()
	r.Mount("/api/customers", customerRouter)

	// License stats endpoint
	r.Get("/api/licenses/stats", customerHandler.LicenseStats)

	// Serve static frontend files (CSS, JS, images, help pages)
	fileServer := http.FileServer(http.Dir(frontendDir))
	r.Handle("/*", fileServer)

	// Serve index.html at root path
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		indexPath := filepath.Join(frontendDir, "index.html")
		http.ServeFile(w, r, indexPath)
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// corsMiddleware adds CORS headers to allow the frontend to access the API.
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
