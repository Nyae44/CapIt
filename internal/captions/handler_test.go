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
		"topic": "wedding photography",
		"tone": "professional"
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
