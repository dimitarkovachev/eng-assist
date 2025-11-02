package ui

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/dimitarkovachev/eng-assist/pkg/state"
	"github.com/gin-gonic/gin"
)

// AssistantUI manages the web interface
type AssistantUI struct {
	router    *gin.Engine
	onPause   func()
	mu        sync.RWMutex
	isRunning bool
	appState  *state.AppState
}

// State represents the current UI state
type State struct {
	Transcript string  `json:"transcript"`
	Response   string  `json:"response"`
	Cost       float64 `json:"cost"`
}

// NewAssistantUI creates a new UI instance
func NewAssistantUI(onPause func(), appState *state.AppState) *AssistantUI {
	ui := &AssistantUI{
		onPause:  onPause,
		appState: appState,
	}

	// Setup Gin router
	router := gin.Default()

	// API routes
	api := router.Group("/api")
	{
		api.GET("/state", ui.getState)
		api.POST("/reset", ui.handleReset)
		api.POST("/pause", ui.handlePause)
	}

	// Serve static files
	router.Static("/static", "ui/static")

	// Serve index.html at root
	router.GET("/", func(c *gin.Context) {
		c.File("ui/static/index.html")
	})

	ui.router = router
	return ui
}

// Run starts the web server
func (ui *AssistantUI) Run(ctx context.Context) error {
	ui.mu.Lock()
	if ui.isRunning {
		ui.mu.Unlock()
		return nil
	}
	ui.isRunning = true
	ui.mu.Unlock()

	server := &http.Server{
		Addr:    ":5001",
		Handler: ui.router,
	}

	// Run server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	fmt.Println("ctx.Done() in Run UI")
	return server.Shutdown(context.Background())
}

// API Handlers

func (ui *AssistantUI) getState(c *gin.Context) {
	ui.mu.RLock()
	defer ui.mu.RUnlock()

	ts, err := ui.appState.TranscriptState.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ars, err := ui.appState.AiResponsesState.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, State{
		Transcript: ts,
		Response:   ars,
		Cost:       0,
	})
}

func (ui *AssistantUI) handleReset(c *gin.Context) {
	ui.appState.Clear()
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (ui *AssistantUI) handlePause(c *gin.Context) {
	ui.appState.SetPaused(!ui.appState.IsPaused())
	if ui.onPause != nil {
		ui.onPause()
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
