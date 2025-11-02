package state

import "sync"

type AppState struct {
	TranscriptState  *TextState
	AiResponsesState *TextState

	mu       sync.RWMutex
	cost     float64
	isPaused bool
}

func NewAppState() *AppState {
	return &AppState{
		TranscriptState:  New(),
		AiResponsesState: New(),
	}
}

func (self *AppState) Clear() {
	self.TranscriptState.Clear()
	self.AiResponsesState.Clear()

	self.mu.Lock()
	defer self.mu.Unlock()

	self.cost = 0
}

func (self *AppState) IsPaused() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.isPaused
}

func (self *AppState) SetPaused(paused bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isPaused = paused
}

func (self *AppState) GetCost() float64 {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.cost
}

func (self *AppState) SetCost(cost float64) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.cost = cost
}
