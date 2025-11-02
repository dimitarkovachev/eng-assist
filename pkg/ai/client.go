package ai

// Tool represents an AI tool interface
type Tool interface {
	// CurrentCost returns the current total cost of API usage
	CurrentCost() float64
}

// MockTool is a mock implementation of the AI Tool interface
type MockTool struct {
}

// NewMockTool creates a new mock AI tool
func NewMockTool() Tool {
	return &MockTool{}
}

// CurrentCost returns the mock cost
func (m *MockTool) CurrentCost() float64 {
	return 0
}
