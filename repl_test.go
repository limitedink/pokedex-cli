package main
import (
	"testing"
	"bytes"
	"os"
	"strings"
)


func TestCleanInput(t *testing.T) {
			cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "HELLO WORLD",
			expected: []string{"hello", "world"},
		},
	}

		for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected){
			t.Errorf("expected length: %d, got %d", len(c.expected), len(actual))
		}
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("expected word: %s, got %s", expectedWord, word)
			}
		// Check each word in the slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		}
	}
}

func TestCommandHelp(t *testing.T) {
	// Save original stdout
	oldStdout := os.Stdout

	// Create a pipe to capture output
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command
	cfg := &config{}
	err := commandHelp(cfg)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "Welcome to the Pokedex!") {
		t.Errorf("help output missing welcome message")
	}

	if !strings.Contains(output, "exit:") {
		t.Errorf("help output missing exit command")
	}

	if !strings.Contains(output, "help:") {
		t.Errorf("help output missing help command")
	}
}

func TestCommandExit(t *testing.T) {
	called := false
	code := -1

	exitFunc = func(c int) {
		called = true
		code = c
	}

	defer func() {
		exitFunc = os.Exit
	}()

	cfg := &config{}
	err := commandExit(cfg)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !called {
		t.Errorf("expected exitFunc to be called")
	}

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}
