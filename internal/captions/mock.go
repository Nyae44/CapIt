package captions

import (
	"fmt"
	"strings"
)

// MockGenerator provides temporary captions for development.
type MockGenerator struct{}

// Generate returns a placeholder caption using the supplied preferences.
func (m MockGenerator) Generate(req GenerateRequest) (string, error) {
	caption := fmt.Sprintf(
		"%s: Discover the story behind %s. "+
			"Experience something special and share the moment!",
		strings.ToUpper(req.Tone),
		req.Topic,
	)

	// Include generation preferences in the mock response.
	metadata := fmt.Sprintf(
		"\n\nPlatform: %s | Content: %s | Language: %s | Length: %s",
		req.Platform,
		req.ContentType,
		req.Language,
		req.CaptionLength,
	)

	return caption + metadata, nil
}
