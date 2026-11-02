package main

import (
	"capit/internal/captions"
	"capit/internal/config"
	"capit/internal/providers"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using environment variables")
	}

	// Load application configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Gemini provider.
	generator, err := providers.NewGeminiGenerator(
		cfg.AI.GeminiAPIKey,
		cfg.AI.Model,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Gemini: %v", err)
	}

	// Initialize router.
	r := chi.NewRouter()

	// Serve static assets.
	assetsDir := http.Dir("./web/assets")
	fileServer := http.FileServer(assetsDir)

	r.Handle("/assets/*", http.StripPrefix("/assets/", fileServer))

	// Serve frontend.
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	// Health check.
	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Initialize caption handler with Gemini.
	captionHandler := captions.NewHandler(generator)

	// Register caption generation endpoint.
	r.Post("/api/v1/captions/generate", captionHandler.Generate)

	// Handle unknown routes.
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// Check that frontend exists.
	indexPath := filepath.Join("web", "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		log.Fatalf("Frontend file not found at %s: %v", indexPath, err)
	}

	// Start server.
	addr := ":" + cfg.Server.Port

	log.Printf("CapIt is running at http://localhost%s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
