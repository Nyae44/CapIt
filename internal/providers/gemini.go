package providers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/genai"
)

// GeminiGenerator generates captions using Google's Gemini API.
type GeminiGenerator struct {
	client *genai.Client
	model  string
}

// NewGeminiGenerator creates a new Gemini provider.
func NewGeminiGenerator(apiKey, model string) (*GeminiGenerator, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("Gemini API key is required")
	}

	if strings.TrimSpace(model) == "" {
		return nil, errors.New("Gemini model is required")
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("create Gemini client: %w", err)
	}

	return &GeminiGenerator{
		client: client,
		model:  model,
	}, nil
}

// Generate creates a social media caption.
func (g *GeminiGenerator) Generate(
	platform, topic, tone string,
) (string, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	prompt := fmt.Sprintf(`
You are an expert social media copywriter.

Create one engaging social media caption using these details:

Platform: %s
Topic: %s
Tone: %s

Requirements:
- Write naturally and creatively.
- Match the requested platform and tone.
- Make the caption engaging and relevant.
- Include an appropriate call to action where suitable.
- Include relevant hashtags where appropriate.
- Return only the caption, without explanations or headings.
`, platform, topic, tone)

	result, err := g.client.Models.GenerateContent(
		ctx,
		g.model,
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("Gemini generation failed: %w", err)
	}

	caption := strings.TrimSpace(result.Text())
	if caption == "" {
		return "", errors.New("Gemini returned an empty caption")
	}

	return caption, nil
}
