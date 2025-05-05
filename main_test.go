package main

import (
	"bytes"
	"os"
	"testing"
)

func createTestFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
		tmpFile.Close()
	})

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	return tmpFile.Name()
}

func captureOutput(t *testing.T, f func()) string {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestSearchFile_BasicMatch(t *testing.T) {
	content := `line1
line2 with pattern
line3
line4 with pattern again`
	filename := createTestFile(t, content)

	output := captureOutput(t, func() {
		searchFile("pattern", filename)
	})

	expected := filename + ":2:line2 with pattern\n" + filename + ":4:line4 with pattern again\n"
	if output != expected {
		t.Errorf("Expected output:\n%q\nGot:\n%q", expected, output)
	}
}

func TestSearchFile_NoMatch(t *testing.T) {
	content := `line1
line2
line3`
	filename := createTestFile(t, content)

	output := captureOutput(t, func() {
		searchFile("pattern", filename)
	})

	if output != "" {
		t.Errorf("Expected no output, got:\n%q", output)
	}
}

func TestSearchFile_EmptyPattern(t *testing.T) {
	content := `line1
line2
line3`
	filename := createTestFile(t, content)

	output := captureOutput(t, func() {
		searchFile("", filename)
	})

	expected := filename + ":1:line1\n" + filename + ":2:line2\n" + filename + ":3:line3\n"
	if output != expected {
		t.Errorf("Expected output:\n%q\nGot:\n%q", expected, output)
	}
}

func TestSearchStdin_BasicMatch(t *testing.T) {
	input := "line1\nline2 with pattern\nline3\nline4 with pattern again\n"

	output := captureOutput(t, func() {
		// Create a pipe for stdin
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("Failed to create pipe: %v", err)
		}

		// Write test input to the pipe
		_, err = w.WriteString(input)
		if err != nil {
			t.Fatalf("Failed to write to pipe: %v", err)
		}
		w.Close()

		// Save and replace stdin
		oldStdin := os.Stdin
		os.Stdin = r
		defer func() { os.Stdin = oldStdin }()

		searchStdin("pattern")
	})

	expected := "stdin:2:line2 with pattern\nstdin:4:line4 with pattern again\n"
	if output != expected {
		t.Errorf("Expected output:\n%q\nGot:\n%q", expected, output)
	}
}

func TestSearchStdin_NoMatch(t *testing.T) {
	input := "line1\nline2\nline3\n"

	output := captureOutput(t, func() {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("Failed to create pipe: %v", err)
		}

		_, err = w.WriteString(input)
		if err != nil {
			t.Fatalf("Failed to write to pipe: %v", err)
		}
		w.Close()

		oldStdin := os.Stdin
		os.Stdin = r
		defer func() { os.Stdin = oldStdin }()

		searchStdin("pattern")
	})

	if output != "" {
		t.Errorf("Expected no output, got:\n%q", output)
	}
}

func TestSearchStdin_EmptyPattern(t *testing.T) {
	input := "line1\nline2\nline3\n"

	output := captureOutput(t, func() {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("Failed to create pipe: %v", err)
		}

		_, err = w.WriteString(input)
		if err != nil {
			t.Fatalf("Failed to write to pipe: %v", err)
		}
		w.Close()

		oldStdin := os.Stdin
		os.Stdin = r
		defer func() { os.Stdin = oldStdin }()

		searchStdin("")
	})

	expected := "stdin:1:line1\nstdin:2:line2\nstdin:3:line3\n"
	if output != expected {
		t.Errorf("Expected output:\n%q\nGot:\n%q", expected, output)
	}
}
