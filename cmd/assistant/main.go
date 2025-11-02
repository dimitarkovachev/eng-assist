package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimitarkovachev/eng-assist/pkg/ai"
	"github.com/dimitarkovachev/eng-assist/pkg/config"
	"github.com/dimitarkovachev/eng-assist/pkg/state"
	"github.com/dimitarkovachev/eng-assist/pkg/transcription"
	"github.com/dimitarkovachev/eng-assist/pkg/ui"
)

// Assistant manages the core application components
type Assistant struct {
	transcription transcription.Transcriptor
	aiClient      ai.Tool
	ui            *ui.AssistantUI
	logger        *log.Logger
	appState      *state.AppState
}

// NewAssistant creates a new assistant instance
func NewAssistant(cfg *config.Config, logger *log.Logger) (*Assistant, error) {
	assistant := &Assistant{
		logger:   logger,
		appState: state.NewAppState(),
	}

	// Initialize AI client with mock implementation
	assistant.aiClient = ai.NewMockTool()

	// Initialize UI with callbacks
	assistant.ui = ui.NewAssistantUI(
		assistant.handlePause,
		assistant.appState,
	)

	// Initialize transcription with callback
	// assistant.transcription = transcription.NewMockTranscriptor(
	// 	strings.Join([]string{
	// 		"hello", "world", "golang", "concurrent", "programming", "state", "management",
	// 		"transcription", "data", "processing", "stream", "buffer", "memory", "async",
	// 		"goroutine", "channel", "mutex", "read", "write", "text", "content", "message",
	// 		"system", "service", "package", "module", "interface", "implementation", "test",
	// 		"example", "demo", "application", "software", "development", "code", "function",
	// 		"method", "variable", "string", "error", "value", "result", "output", "input",
	// 	}, " "),
	// 	assistant.appState,
	// )

	assistant.transcription = transcription.NewWhisperHandler(
		cfg.WhisperCppPath,
		cfg.WhisperModelPath,
		cfg.BufferTimeout,
		assistant.appState,
	)

	return assistant, nil
}

// handlePause toggles the pause state
func (a *Assistant) handlePause() {
	// if a.appState.IsPaused() {
	// 	a.transcription.Stop()
	// } else {
	// 	ctx := context.Background()
	// 	if err := a.transcription.Start(ctx); err != nil {
	// 		a.logger.Printf("Error starting transcription: %v", err.Error())
	// 	}
	// }
}

// Run starts the assistant
func (a *Assistant) Run(ctx context.Context) error {
	// Start transcription
	if err := a.transcription.Start(ctx); err != nil {
		return err
	}

	// Start UI
	if err := a.ui.Run(ctx); err != nil {
		return err
	}

	return nil
}

func main() {
	// Parse flags
	bufferTimeout := flag.Float64("buffer-timeout", 1.0, "Time to wait before processing buffered text")
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	// Setup logging
	logFlags := log.LstdFlags
	if *debug {
		logFlags |= log.Lshortfile
	}
	logger := log.New(os.Stdout, "", logFlags)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}
	cfg.BufferTimeout = *bufferTimeout
	cfg.Debug = *debug

	logger.Printf("Config: %+v", cfg)

	// Create and start assistant
	assistant, err := NewAssistant(cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to create assistant: %v", err)
	}

	// Setup context with signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
		fmt.Println("sigChan done")
	}()

	// Run assistant
	if err := assistant.Run(ctx); err != nil && err != context.Canceled {
		logger.Fatalf("Assistant error: %v", err)
	}

	fmt.Println("Waiting 5s for graceful shutdown")
	time.Sleep(5 * time.Second)
	fmt.Println("main done")
}
