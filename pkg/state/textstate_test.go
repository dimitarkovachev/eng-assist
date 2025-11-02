package state

import (
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ts := New()
	if ts == nil {
		t.Fatal("New() returned nil")
	}

	state, hasNew, err := ts.Read()
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}
	if state != "" {
		t.Errorf("Expected empty state, got: %s", state)
	}
	if hasNew {
		t.Error("Expected hasNew to be false for new instance")
	}
}

func TestWrite(t *testing.T) {
	ts := New()

	err := ts.Write("hello")
	if err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	state, hasNew, err := ts.Read()
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}
	if state != "hello" {
		t.Errorf("Expected 'hello', got: %s", state)
	}
	if !hasNew {
		t.Error("Expected hasNew to be true after write")
	}
}

func TestWriteAppend(t *testing.T) {
	ts := New()

	ts.Write("hello")
	ts.Write(" world")

	state, hasNew, err := ts.Read()
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}
	if state != "hello world" {
		t.Errorf("Expected 'hello world', got: %s", state)
	}
	if !hasNew {
		t.Error("Expected hasNew to be true after write")
	}
}

func TestReadResetsHasNew(t *testing.T) {
	ts := New()

	ts.Write("test")

	state1, hasNew1, err := ts.Read()
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}
	if !hasNew1 {
		t.Error("Expected hasNew to be true after write")
	}

	state2, hasNew2, err := ts.Read()
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}
	if hasNew2 {
		t.Error("Expected hasNew to be false after second read")
	}
	if state1 != state2 {
		t.Errorf("States should be equal: %s vs %s", state1, state2)
	}
}

func TestTrimToMaxWords(t *testing.T) {
	ts := New()

	words := make([]string, 250)
	for i := 0; i < 250; i++ {
		words[i] = "word"
	}
	longText := strings.Join(words, " ")

	ts.Write(longText)

	state, _, err := ts.Read()
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	resultWords := strings.Fields(state)
	if len(resultWords) != 200 {
		t.Errorf("Expected 200 words, got: %d", len(resultWords))
	}
}

func TestTrimKeepsLatestWords(t *testing.T) {
	ts := New()

	ts.Write("old words that should be removed ")

	words := make([]string, 200)
	for i := 0; i < 200; i++ {
		words[i] = "new"
	}
	ts.Write(strings.Join(words, " "))

	state, _, err := ts.Read()
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	if strings.Contains(state, "old") {
		t.Error("Old words should have been trimmed")
	}

	resultWords := strings.Fields(state)
	if len(resultWords) != 200 {
		t.Errorf("Expected 200 words, got: %d", len(resultWords))
	}
}

func TestEmptyWrite(t *testing.T) {
	ts := New()

	err := ts.Write("")
	if err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	state, hasNew, err := ts.Read()
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}
	if state != "" {
		t.Errorf("Expected empty state, got: %s", state)
	}
	if !hasNew {
		t.Error("Expected hasNew to be true even for empty write")
	}
}

func TestConcurrentWrites(t *testing.T) {
	ts := New()

	var wg sync.WaitGroup
	numGoroutines := 10
	writesPerGoroutine := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < writesPerGoroutine; j++ {
				ts.Write("test ")
			}
		}(i)
	}

	wg.Wait()

	state, hasNew, err := ts.Read()
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}
	if !hasNew {
		t.Error("Expected hasNew to be true after concurrent writes")
	}

	words := strings.Fields(state)
	expectedWords := numGoroutines * writesPerGoroutine
	if expectedWords > 200 {
		expectedWords = 200
	}

	if len(words) != expectedWords {
		t.Errorf("Expected %d words, got: %d", expectedWords, len(words))
	}
}

func TestConcurrentReadsAndWrites(t *testing.T) {
	ts := New()

	var wg sync.WaitGroup
	done := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			ts.Write("concurrent ")
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				ts.Read()
				time.Sleep(time.Millisecond)
			}
		}
	}()

	wg.Wait()

	state, _, err := ts.Read()
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	words := strings.Fields(state)
	if len(words) == 0 {
		t.Error("Expected some words after concurrent operations")
	}
}

func TestWordBoundaryTrimming(t *testing.T) {
	ts := New()

	text := "This is a very long sentence that should be trimmed properly at word boundaries. "
	for i := 0; i < 30; i++ {
		ts.Write(text)
	}

	state, _, err := ts.Read()
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	words := strings.Fields(state)
	if len(words) != 200 {
		t.Errorf("Expected exactly 200 words, got: %d", len(words))
	}

	if !strings.HasSuffix(state, "boundaries.") {
		t.Error("Expected trimming to preserve word boundaries")
	}
}

func TestWhitespaceHandling(t *testing.T) {
	ts := New()

	ts.Write("  hello   world  ")

	state, _, err := ts.Read()
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	if state != "  hello   world  " {
		t.Errorf("Expected whitespace to be preserved, got: '%s'", state)
	}
}

func TestMultipleEmptyReads(t *testing.T) {
	ts := New()

	for i := 0; i < 5; i++ {
		state, hasNew, err := ts.Read()
		if err != nil {
			t.Fatalf("Read() error: %v", err)
		}
		if state != "" {
			t.Errorf("Expected empty state, got: %s", state)
		}
		if hasNew {
			t.Error("Expected hasNew to be false for empty reads")
		}
	}
}
