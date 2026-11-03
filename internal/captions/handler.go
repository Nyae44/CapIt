package captions

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

// GenerateRequest represents the incoming caption generation request.
type GenerateRequest struct {
	Platform      string `json:"platform"`
	ContentType   string `json:"content_type"`
	Topic         string `json:"topic"`
	Tone          string `json:"tone"`
	Language      string `json:"language"`
	CaptionLength string `json:"caption_length"`
}

// GenerateResponse represents the API response.
type GenerateResponse struct {
	Caption  string `json:"caption"`
	Platform string `json:"platform"`
}

// Generator defines how captions are generated.
type Generator interface {
	Generate(req GenerateRequest) (string, error)
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

	// Decode the incoming JSON request.
	var req GenerateRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON request")
		return
	}

	// Ensure the request contains only one JSON value.
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeJSONError(w, http.StatusBadRequest, "request must contain a single JSON object")
		return
	}

	// Validate the request before calling the generator.
	if err := ValidateRequest(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Generate the caption.
	caption, err := h.generator.Generate(req)
	if err != nil {
		log.Printf("Caption generation failed: %v", err)

		writeJSONError(
			w,
			http.StatusInternalServerError,
			"failed to generate caption",
		)
		return
	}

	// Return the generated caption.
	writeJSON(w, http.StatusOK, GenerateResponse{
		Caption:  caption,
		Platform: req.Platform,
	})
}

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON response encoding failed: %v", err)
	}
}

// writeJSONError writes a standardized JSON error response.
func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}
