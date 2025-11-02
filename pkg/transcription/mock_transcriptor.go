package transcription

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dimitarkovachev/eng-assist/pkg/state"
)

// MockTranscriptor simulates a transcription process by outputting text at a natural speaking rate
type MockTranscriptor struct {
	text     string
	appState *state.AppState

	mu        sync.Mutex
	isRunning bool
	cancel    context.CancelFunc
}

// NewMockTranscriptor creates a new mock transcriptor
func NewMockTranscriptor(text string, appState *state.AppState) *MockTranscriptor {
	return &MockTranscriptor{
		text:     text,
		appState: appState,
	}
}

// Start begins the mock transcription process
func (m *MockTranscriptor) Start(ctx context.Context) error {
	m.mu.Lock()
	if m.isRunning {
		m.mu.Unlock()
		return nil
	}
	m.isRunning = true
	m.mu.Unlock()

	ctx, m.cancel = context.WithCancel(ctx)

	go m.simulateSpeaking(ctx)

	return nil
}

// Stop terminates the mock transcription process
func (m *MockTranscriptor) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return nil
	}

	if m.cancel != nil {
		m.cancel()
	}

	m.isRunning = false
	return nil
}

func (m *MockTranscriptor) simulateSpeaking(ctx context.Context) {
	// Average speaking rate is 150 words per minute
	// That's 2.5 words per second
	// const wordsPerSecond = 2.5

	words := strings.Fields(m.text)
	currentWord := 0

	ticker := time.NewTicker(time.Duration(1000000000))
	defer ticker.Stop()

	// 3 words every 1 second

	for {
		select {
		case <-ctx.Done():
			fmt.Println("ctx.Done() in simulateSpeaking")
			return
		case <-ticker.C:
			if currentWord+3 > len(words) {
				currentWord = 0
			}

			// Write the current buffer to the transcript state
			m.appState.TranscriptState.Write(strings.Join(words[currentWord:currentWord+3], " "))

			currentWord += 3
		}
	}
}
