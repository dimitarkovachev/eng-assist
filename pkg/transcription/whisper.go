package transcription

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/creack/pty"
	"github.com/dimitarkovachev/eng-assist/pkg/state"
)

// WhisperHandler manages the whisper.cpp transcription process
type WhisperHandler struct {
	whisperPath   string
	modelPath     string
	bufferTimeout time.Duration
	cmd           *exec.Cmd
	appState      *state.AppState

	mu        sync.Mutex
	cancel    context.CancelFunc
	isRunning bool
}

// NewWhisperHandler creates a new transcription handler
func NewWhisperHandler(whisperPath string, modelPath string, bufferTimeout float64, appState *state.AppState) *WhisperHandler {
	return &WhisperHandler{
		whisperPath:   whisperPath,
		modelPath:     modelPath,
		bufferTimeout: time.Duration(bufferTimeout * float64(time.Second)),
		appState:      appState,
	}
}

// Start begins the whisper.cpp transcription process
func (h *WhisperHandler) Start(ctx context.Context) error {
	h.mu.Lock()
	if h.isRunning {
		h.mu.Unlock()
		return nil
	}
	h.isRunning = true
	h.mu.Unlock()

	deviceID := GetDeviceIDByName("VB-Cable")

	ctx, h.cancel = context.WithCancel(ctx)
	h.cmd = exec.CommandContext(ctx, h.whisperPath,
		// "-t", "8",
		"-m", h.modelPath,
		"-c", strconv.Itoa(deviceID),
	)

	ptmx, err := pty.Start(h.cmd)
	if err != nil {
		return fmt.Errorf("failed to start pty: %w", err)
	}

	// stdoutPipe, err := h.cmd.StdoutPipe()
	// if err != nil {
	// 	return fmt.Errorf("failed to get stdout pipe: %w", err)
	// }
	// stderrPipe, err := h.cmd.StderrPipe()
	// if err != nil {
	// 	return fmt.Errorf("failed to get stderr pipe: %w", err)
	// }

	// if err := h.cmd.Start(); err != nil {
	// 	return err
	// }

	// Goroutine to read stdout and update app state
	go func() {
		scanner := bufio.NewScanner(ptmx)
		for scanner.Scan() {
			line := scanner.Text()

			// // Handle clear-line + carriage return (\x1b[2K\r)
			// if strings.Contains(line, "\x1b[2K") {
			// 	h.appState.TranscriptState.RemoveLastLine()
			// 	clean := strings.ReplaceAll(line, "\x1b[2K", "")
			// 	clean = strings.TrimLeft(clean, "\r")
			// 	if clean != "" {
			// 		h.appState.TranscriptState.Write(clean)
			// 	}
			// 	continue
			// }

			// // If the line starts with a carriage return, treat as overwrite (delete previous line)
			// if strings.HasPrefix(line, "\r") {
			// 	h.appState.TranscriptState.RemoveLastLine()
			// 	clean := strings.TrimLeft(line, "\r")
			// 	if clean != "" {
			// 		h.appState.TranscriptState.Write(clean)
			// 	}
			// 	continue
			// }

			h.appState.TranscriptState.Write(line)
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading stdout: %v\n", err)
		}
	}()

	// Goroutine to read and log stderr
	// go func() {
	// 	scanner := bufio.NewScanner(stderrPipe)
	// 	for scanner.Scan() {
	// 		line := scanner.Text()
	// 		fmt.Printf("[whisper stderr] %s\n", line)
	// 	}
	// 	if err := scanner.Err(); err != nil {
	// 		fmt.Printf("Error reading stderr: %v\n", err)
	// 	}
	// }()

	return nil
}

// Stop terminates the transcription process
func (h *WhisperHandler) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.isRunning {
		return nil
	}

	if h.cancel != nil {
		h.cancel()
	}

	if h.cmd != nil && h.cmd.Process != nil {
		if err := h.cmd.Process.Kill(); err != nil {
			return err
		}
		h.cmd.Wait()
	}

	h.isRunning = false
	return nil
}
