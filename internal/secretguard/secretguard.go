package secretguard

import (
	"regexp"
	"strings"
)

// SecretType represents the type of detected secret
type SecretType string

const (
	SecretTypeAWS       SecretType = "AWS Access Key"
	SecretTypeGitHub    SecretType = "GitHub Token"
	SecretTypeAnthropic SecretType = "Anthropic API Key"
	SecretTypeOpenAI    SecretType = "OpenAI API Key"
	SecretTypeGeneric   SecretType = "Generic API Key"
	SecretTypePrivate   SecretType = "Private Key"
	SecretTypeJWT       SecretType = "JWT Token"
)

// Detection represents a detected secret
type Detection struct {
	Type     SecretType
	Original string
	Masked   string
	Start    int
	End      int
}

// patterns holds compiled regex patterns for secret detection
var patterns = map[SecretType]*regexp.Regexp{
	// AWS Access Key ID
	SecretTypeAWS: regexp.MustCompile(`AKIA[0-9A-Z]{16}`),

	// GitHub Tokens (classic and fine-grained)
	SecretTypeGitHub: regexp.MustCompile(`gh[pousr]_[A-Za-z0-9]{36,}`),

	// Anthropic API Keys
	SecretTypeAnthropic: regexp.MustCompile(`sk-ant-[A-Za-z0-9-]{20,}`),

	// OpenAI API Keys
	SecretTypeOpenAI: regexp.MustCompile(`sk-[A-Za-z0-9]{32,}`),

	// Generic API Key patterns
	SecretTypeGeneric: regexp.MustCompile(`(?i)(api[_-]?key|apikey|api_secret|access_token)[=:]["']?([A-Za-z0-9_-]{20,})["']?`),

	// Private Keys (PEM format)
	SecretTypePrivate: regexp.MustCompile(`-----BEGIN[A-Z ]+PRIVATE KEY-----`),

	// JWT Tokens
	SecretTypeJWT: regexp.MustCompile(`eyJ[A-Za-z0-9_-]*\.eyJ[A-Za-z0-9_-]*\.[A-Za-z0-9_-]*`),
}

// Guard handles secret detection and masking
type Guard struct {
	enabled bool
}

// New creates a new SecretGuard
func New(enabled bool) *Guard {
	return &Guard{enabled: enabled}
}

// Scan checks text for secrets and returns detections
func (g *Guard) Scan(text string) []Detection {
	if !g.enabled {
		return nil
	}

	var detections []Detection

	for secretType, pattern := range patterns {
		matches := pattern.FindAllStringIndex(text, -1)
		for _, match := range matches {
			original := text[match[0]:match[1]]
			detections = append(detections, Detection{
				Type:     secretType,
				Original: original,
				Masked:   maskSecret(original),
				Start:    match[0],
				End:      match[1],
			})
		}
	}

	return detections
}

// Mask replaces all detected secrets in text with masked versions
func (g *Guard) Mask(text string) string {
	if !g.enabled {
		return text
	}

	detections := g.Scan(text)

	// Sort by position descending to replace from end to start
	// (avoids index shifting issues)
	for i := len(detections) - 1; i >= 0; i-- {
		d := detections[i]
		text = text[:d.Start] + d.Masked + text[d.End:]
	}

	return text
}

// maskSecret creates a masked version of a secret
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "****"
	}

	// Show first 4 chars, mask middle, show last 4
	prefix := secret[:4]
	suffix := secret[len(secret)-4:]
	middle := strings.Repeat("*", min(len(secret)-8, 20))

	return prefix + middle + suffix
}

// HasSecrets returns true if the text contains any secrets
func (g *Guard) HasSecrets(text string) bool {
	return len(g.Scan(text)) > 0
}

// Enable enables secret detection
func (g *Guard) Enable() {
	g.enabled = true
}

// Disable disables secret detection
func (g *Guard) Disable() {
	g.enabled = false
}

// IsEnabled returns the current state
func (g *Guard) IsEnabled() bool {
	return g.enabled
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
