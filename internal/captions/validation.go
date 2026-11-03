package captions

import (
	"fmt"
	"strings"
)

// ValidateRequest validates and normalizes caption generation input.
func ValidateRequest(req *GenerateRequest) error {
	req.Platform = strings.ToLower(strings.TrimSpace(req.Platform))
	req.ContentType = strings.ToLower(strings.TrimSpace(req.ContentType))
	req.Topic = strings.TrimSpace(req.Topic)
	req.Tone = strings.ToLower(strings.TrimSpace(req.Tone))
	req.Language = strings.ToLower(strings.TrimSpace(req.Language))
	req.CaptionLength = strings.ToLower(strings.TrimSpace(req.CaptionLength))

	if req.Topic == "" {
		return fmt.Errorf("topic is required")
	}

	if len(req.Topic) > 500 {
		return fmt.Errorf("topic must not exceed 500 characters")
	}

	if !isValid(req.Platform, []string{
		"instagram", "facebook", "tiktok", "linkedin",
	}) {
		return fmt.Errorf("unsupported platform")
	}

	if !isValid(req.ContentType, []string{
		"photo", "video", "carousel", "reel", "short", "story", "text",
	}) {
		return fmt.Errorf("unsupported content type")
	}

	if !isValid(req.Tone, []string{
		"friendly", "professional", "playful",
		"inspirational", "witty", "bold",
	}) {
		return fmt.Errorf("unsupported tone")
	}

	if !isValid(req.Language, []string{
		"english", "german", "french",
		"portuguese", "spanish",
	}) {
		return fmt.Errorf("unsupported language")
	}

	if !isValid(req.CaptionLength, []string{
		"short", "medium", "long",
	}) {
		return fmt.Errorf("unsupported caption length")
	}

	return nil
}

// isValid checks whether a value exists in the allowed options.
func isValid(value string, allowed []string) bool {
	for _, option := range allowed {
		if value == option {
			return true
		}
	}

	return false
}
