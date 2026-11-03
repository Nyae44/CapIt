package captions

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenerateCaption(t *testing.T) {
	handler := NewHandler(MockGenerator{})

	body := `{
		"platform": "instagram",
		"content_type": "photo",
		"topic": "wedding photography",
		"tone": "professional",
		"language": "english",
		"caption_length": "medium"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/captions/generate",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.Generate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s",
			rec.Code, rec.Body.String())
	}

	var response GenerateResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Caption == "" {
		t.Error("expected a generated caption, got an empty string")
	}

	if response.Platform != "instagram" {
		t.Errorf("expected platform instagram, got %q",
			response.Platform)
	}
}

func TestGenerateCaptionMissingTopic(t *testing.T) {
	handler := NewHandler(MockGenerator{})

	body := `{
		"platform": "instagram",
		"content_type": "photo",
		"topic": "",
		"tone": "professional",
		"language": "english",
		"caption_length": "medium"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/captions/generate",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.Generate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestGenerateCaptionUnsupportedPlatform(t *testing.T) {
	handler := NewHandler(MockGenerator{})

	body := `{
		"platform": "twitter",
		"content_type": "photo",
		"topic": "wedding photography",
		"tone": "professional",
		"language": "english",
		"caption_length": "medium"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/captions/generate",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.Generate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestGenerateCaptionInvalidJSON(t *testing.T) {
	handler := NewHandler(MockGenerator{})

	body := `{"platform": "instagram",`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/captions/generate",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.Generate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestGenerateCaptionUnsupportedContentType(t *testing.T) {
	handler := NewHandler(MockGenerator{})

	body := `{
		"platform": "instagram",
		"content_type": "document",
		"topic": "wedding photography",
		"tone": "professional",
		"language": "english",
		"caption_length": "medium"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/captions/generate",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.Generate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}
