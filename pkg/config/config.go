package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration values
type Config struct {
	WhisperCppPath   string
	WhisperModelPath string
	AnthropicApiKey  string // Optional: only needed when using real AI client
	BufferTimeout    float64
	Debug            bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	whisperPath := os.Getenv("WHISPER_CPP_PATH")
	if whisperPath == "" {
		return nil, fmt.Errorf("WHISPER_CPP_PATH environment variable not set")
	}

	modelPath := os.Getenv("WHISPER_MODEL_PATH")
	if modelPath == "" {
		return nil, fmt.Errorf("WHISPER_MODEL_PATH environment variable not set")
	}

	// API key is now optional
	apiKey := os.Getenv("ANTHROPIC_API_KEY")

	return &Config{
		WhisperCppPath:   whisperPath,
		WhisperModelPath: modelPath,
		AnthropicApiKey:  apiKey,
		BufferTimeout:    1.0, // Default value
		Debug:            false,
	}, nil
}
