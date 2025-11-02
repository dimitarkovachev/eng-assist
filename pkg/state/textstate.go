package state

import (
	"strings"
	"sync"
)

type TextState struct {
	mu         sync.RWMutex
	state      string
	hasNewData bool
}

func New() *TextState {
	return &TextState{
		hasNewData: false,
	}
}

func (ts *TextState) Write(txt string) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.state += txt
	ts.hasNewData = true

	ts.trimToMaxWords(500)

	return nil
}

func (ts *TextState) Read() (string, bool, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	state := ts.state
	hasNew := ts.hasNewData
	ts.hasNewData = false

	return state, hasNew, nil
}

func (ts *TextState) trimToMaxWords(maxWords int) {
	words := strings.Fields(ts.state)
	if len(words) > maxWords {
		ts.state = strings.Join(words[len(words)-maxWords:], " ")
	}
}

func (ts *TextState) GetAll() (string, error) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	return ts.state, nil
}

func (ts *TextState) Clear() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.state = ""
	ts.hasNewData = false
}

func (ts *TextState) RemoveLastLine() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	lines := strings.Split(ts.state, "\n")
	if len(lines) > 1 {
		// Remove the last non-empty line
		lines = lines[:len(lines)-2] // Remove last line and the empty string after last \n
		ts.state = strings.Join(lines, "\n")
		if ts.state != "" {
			ts.state += "\n"
		}
	} else {
		ts.state = ""
	}
	// Mark as new data for UI update
	ts.hasNewData = true
}
