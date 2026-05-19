package wizard

import (
	"io"
	"os"
	"testing"
)

func TestWizardValidInput(t *testing.T) {
	input := "example.com\nuser@example.com\n1\n"

	r, w, _ := os.Pipe()
	w.WriteString(input)
	w.Close()

	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	// Suppress stdout
	oldStdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = oldStdout }()

	answers, err := Run()
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if answers.Domain != "example.com" {
		t.Errorf("Domain: got %q, want %q", answers.Domain, "example.com")
	}
	if answers.Email != "user@example.com" {
		t.Errorf("Email: got %q, want %q", answers.Email, "user@example.com")
	}
	if answers.Challenge != "http-01" {
		t.Errorf("Challenge: got %q, want %q", answers.Challenge, "http-01")
	}
}

func TestWizardDNSChoice(t *testing.T) {
	input := "example.com\nuser@example.com\n2\n"

	r, w, _ := os.Pipe()
	w.WriteString(input)
	w.Close()

	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	oldStdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = oldStdout }()

	answers, err := Run()
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if answers.Challenge != "dns-01" {
		t.Errorf("Challenge: got %q, want %q", answers.Challenge, "dns-01")
	}
}

func TestWizardEmptyDomain(t *testing.T) {
	input := "\n"

	r, w, _ := os.Pipe()
	w.WriteString(input)
	w.Close()

	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	oldStdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = oldStdout }()

	_, err := Run()
	if err == nil {
		t.Fatal("Expected error for empty domain")
	}
}

func TestWizardEOFInput(t *testing.T) {
	r, w, _ := os.Pipe()
	w.Close()

	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	oldStdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = oldStdout }()

	_, err := Run()
	if err == nil {
		t.Fatal("Expected error for EOF")
	}
	_ = io.Discard
}
