package main

import (
	"capit/internal/captions"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// Serve static assets.
	assetsDir := http.Dir("./web/assets")
	fileServer := http.FileServer(assetsDir)

	r.Handle("/assets/*", http.StripPrefix("/assets/", fileServer))

	// Serve the frontend.
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	// Health-check endpoint.
	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Initalize the caption generator
	generator := captions.MockGenerator{}

	// Initialize the caption handler.
	captionHandler := captions.NewHandler(generator)

	// Register the caption generation endpoint.
	r.Post("/api/v1/captions/generate", captionHandler.Generate)

	// Handle unknown routes.
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// Check that the frontend exists before starting.
	indexPath := filepath.Join("web", "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		log.Fatalf("Frontend file not found at %s: %v", indexPath, err)
	}

	log.Println("CapIt is running at http://localhost:3000")

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
