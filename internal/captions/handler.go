package captions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GenerateRequest represents the incoming JSON request.
type GenerateRequest struct {
	Platform string `json:"platform"`
	Topic    string `json:"topic"`
	Tone     string `json:"tone"`
}

// GenerateResponse represents the API response.
type GenerateResponse struct {
	Caption  string `json:"caption"`
	Platform string `json:"platform"`
}

// Generator defines how captions are generated.
type Generator interface {
	Generate(platform, topic, tone string) (string, error)
}

// Handler handles caption-related HTTP requests.
type Handler struct {
	generator Generator
}

// NewHandler creates a caption handler.
func NewHandler(generator Generator) *Handler {
	return &Handler{
		generator: generator,
	}
}

// Generate handles POST /api/v1/captions/generate.
func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {
	// Limit the request body to 1 MB.
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	// Decode the incoming JSON.
	var req GenerateRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON request"}`,
			http.StatusBadRequest)
		return
	}

	// Validate required fields.
	req.Platform = strings.TrimSpace(req.Platform)
	req.Topic = strings.TrimSpace(req.Topic)
	req.Tone = strings.TrimSpace(req.Tone)

	if req.Platform == "" || req.Topic == "" || req.Tone == "" {
		http.Error(w,
			`{"error":"platform, topic and tone are required"}`,
			http.StatusBadRequest,
		)
		return
	}

	if len(req.Topic) > 500 {
		http.Error(w,
			`{"error":"topic must not exceed 500 characters"}`,
			http.StatusBadRequest,
		)
		return
	}

	// Validate the supported platforms.
	switch strings.ToLower(req.Platform) {
	case "instagram", "facebook", "tiktok", "linkedin":
		req.Platform = strings.ToLower(req.Platform)
	default:
		http.Error(w,
			`{"error":"unsupported platform"}`,
			http.StatusBadRequest,
		)
		return
	}

	// Generate the caption.
	caption, err := h.generator.Generate(
		req.Platform,
		req.Topic,
		req.Tone,
	)

	if err != nil {
		http.Error(w,
			`{"error":"failed to generate caption"}`,
			http.StatusInternalServerError,
		)
		return
	}

	// Return the response as JSON.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(GenerateResponse{
		Caption:  caption,
		Platform: req.Platform,
	}); err != nil {
		fmt.Printf("encode response: %v\n", err)
	}
}
