package providers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"capit/internal/captions"

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
	req captions.GenerateRequest,
) (string, error) {

	prompt := fmt.Sprintf(`
You are CapIt, an expert social media copywriter.

Create one engaging social media caption using these details:

Platform: %s
Content Type: %s
Topic: %s
Tone: %s
Language: %s
Caption Length: %s

Requirements:
- Adapt the caption to the selected social media platform.
- Respect the selected content type.
- Match the requested tone.
- Write entirely in the requested language.
- Respect the requested caption length.
- Make the caption engaging, natural and relevant.
- Include an appropriate call to action where suitable.
- Include relevant hashtags where appropriate.
- Avoid generic or repetitive expressions.
- Return only the caption, without explanations or headings.
`,
		req.Platform,
		req.ContentType,
		req.Topic,
		req.Tone,
		req.Language,
		req.CaptionLength,
	)

	const maxAttempts = 3

	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			60*time.Second,
		)

		result, err := g.client.Models.GenerateContent(
			ctx,
			g.model,
			genai.Text(prompt),
			nil,
		)

		cancel()

		if err == nil {
			caption := strings.TrimSpace(result.Text())

			if caption == "" {
				return "", errors.New("Gemini returned an empty caption")
			}

			return caption, nil
		}

		lastErr = err

		if !isRetryableError(err) {
			break
		}

		if attempt < maxAttempts {
			delay := time.Duration(1<<uint(attempt-1)) * time.Second

			time.Sleep(delay)
		}
	}

	return "", fmt.Errorf(
		"Gemini generation failed after %d attempt(s): %w",
		maxAttempts,
		lastErr,
	)
}

// isRetryableError identifies temporary provider failures.
func isRetryableError(err error) bool {
	message := strings.ToLower(err.Error())

	retryableErrors := []string{
		"503",
		"429",
		"unavailable",
		"resource_exhausted",
		"deadline exceeded",
		"timeout",
		"connection reset",
		"connection refused",
		"temporarily unavailable",
	}

	for _, retryable := range retryableErrors {
		if strings.Contains(message, retryable) {
			return true
		}
	}

	return false
}
