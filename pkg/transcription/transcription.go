package transcription

import (
	"context"
)

// Transcriptor defines the interface for transcription handlers
type Transcriptor interface {
	// Start begins the transcription process
	Start(ctx context.Context) error

	// Stop terminates the transcription process
	Stop() error
}
