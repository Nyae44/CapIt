package captions

import (
	"fmt"
	"strings"
)

// MockGenerator provides temporary captions for development.
type MockGenerator struct{}

// Generate returns a placeholder caption.
func (m MockGenerator) Generate(
	platform, topic, tone string,
) (string, error) {
	caption := fmt.Sprintf(
		"%s: Discover the story behind %s. "+
			"Experience something special and share the moment!",
		strings.ToUpper(tone),
		topic,
	)

	return caption, nil
}
