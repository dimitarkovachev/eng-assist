package transcription

// This file shows how to integrate whisper.cpp directly using CGO
// To use this, you need to:
// 1. Install whisper.cpp on your system (see WHISPER_INTEGRATION.md)
// 2. Uncomment the CGO directives below
// 3. Build with CGO_ENABLED=1

/*
Example CGO integration (uncomment when whisper.cpp is installed):

#cgo CXXFLAGS: -std=c++11
#cgo darwin LDFLAGS: -lwhisper
#cgo linux LDFLAGS: -lwhisper
#cgo windows LDFLAGS: -lwhisper
#cgo CFLAGS: -I/usr/local/include
#cgo CXXFLAGS: -I/usr/local/include
#include <whisper.h>
#include <stdlib.h>
#include <stdbool.h>
import "C"

// WhisperContext represents a whisper context
type WhisperContext struct {
	ctx *C.struct_whisper_context
}

// NewWhisperContext creates a new whisper context from a model file
func NewWhisperContext(modelPath string) (*WhisperContext, error) {
	cPath := C.CString(modelPath)
	defer C.free(unsafe.Pointer(cPath))

	ctx := C.whisper_init_from_file(cPath)
	if ctx == nil {
		return nil, errors.New("failed to initialize whisper context")
	}

	return &WhisperContext{ctx: ctx}, nil
}

// ProcessAudio processes audio samples and returns transcription
func (w *WhisperContext) ProcessAudio(samples []float32, language string, translate bool) error {
	if w.ctx == nil {
		return errors.New("whisper context is nil")
	}

	// Convert Go slice to C array
	cSamples := (*C.float)(unsafe.Pointer(&samples[0]))

	// Create parameters
	params := C.whisper_full_params_default()

	if language != "" {
		cLang := C.CString(language)
		defer C.free(unsafe.Pointer(cLang))
		params.language = cLang
	}

	params.translate = C.bool(translate)

	// Process audio
	result := C.whisper_full(w.ctx, params, cSamples, C.int(len(samples)))
	if result != 0 {
		return errors.New("failed to process audio")
	}

	return nil
}

// GetTranscription returns the transcription text for a specific segment
func (w *WhisperContext) GetTranscription(segmentID int) (string, error) {
	if w.ctx == nil {
		return "", errors.New("whisper context is nil")
	}

	cText := C.whisper_full_get_text(w.ctx, C.int(segmentID))
	if cText == nil {
		return "", errors.New("failed to get transcription text")
	}

	return C.GoString(cText), nil
}

// GetSegmentCount returns the number of transcription segments
func (w *WhisperContext) GetSegmentCount() (int, error) {
	if w.ctx == nil {
		return 0, errors.New("whisper context is nil")
	}

	count := C.whisper_full_get_n_segments(w.ctx)
	return int(count), nil
}

// GetAllTranscriptions returns all transcription segments
func (w *WhisperContext) GetAllTranscriptions() ([]string, error) {
	count, err := w.GetSegmentCount()
	if err != nil {
		return nil, err
	}

	transcriptions := make([]string, count)
	for i := 0; i < count; i++ {
		text, err := w.GetTranscription(i)
		if err != nil {
			return nil, err
		}
		transcriptions[i] = text
	}

	return transcriptions, nil
}

// Free releases the whisper context
func (w *WhisperContext) Free() {
	if w.ctx != nil {
		C.whisper_free(w.ctx)
		w.ctx = nil
	}
}
*/

// Example usage in WhisperHandler:
/*
func (h *WhisperHandler) StartDirect(ctx context.Context) error {
	h.mu.Lock()
	if h.isRunning {
		h.mu.Unlock()
		return nil
	}
	h.isRunning = true
	h.mu.Unlock()

	// Initialize whisper context directly
	whisperCtx, err := NewWhisperContext(h.modelPath)
	if err != nil {
		return fmt.Errorf("failed to initialize whisper context: %w", err)
	}
	h.whisperCtx = whisperCtx

	// Start audio processing goroutine
	go h.processAudioStream()

	return nil
}

func (h *WhisperHandler) processAudioStream() {
	// This would be called with audio chunks from your audio input
	// For example, when you receive audio from VB-Cable

	// Process audio chunk
	err := h.whisperCtx.ProcessAudio(audioChunk, "en", false)
	if err != nil {
		fmt.Printf("Error processing audio: %v\n", err)
		return
	}

	// Get transcriptions
	transcriptions, err := h.whisperCtx.GetAllTranscriptions()
	if err != nil {
		fmt.Printf("Error getting transcriptions: %v\n", err)
		return
	}

	// Update app state
	for _, text := range transcriptions {
		h.appState.TranscriptState.Write(text)
	}
}

func (h *WhisperHandler) StopDirect() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.isRunning {
		return nil
	}

	if h.whisperCtx != nil {
		h.whisperCtx.Free()
		h.whisperCtx = nil
	}

	h.isRunning = false
	return nil
}
*/
