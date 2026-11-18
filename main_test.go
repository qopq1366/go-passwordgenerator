package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunDefaults(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run(nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() code = %d, want 0; stderr = %q", code, stderr.String())
	}
	generated := strings.TrimSpace(stdout.String())
	if len(generated) != 16 {
		t.Fatalf("generated password length = %d, want 16", len(generated))
	}
}

func TestRunSettings(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"--length", "24",
		"--count", "3",
		"--lowercase=false",
		"--uppercase=false",
		"--symbols=false",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() code = %d, want 0; stderr = %q", code, stderr.String())
	}

	lines := strings.Fields(stdout.String())
	if len(lines) != 3 {
		t.Fatalf("generated password count = %d, want 3", len(lines))
	}
	for _, line := range lines {
		if len(line) != 24 {
			t.Errorf("password length = %d, want 24", len(line))
		}
		for _, character := range line {
			if character < '0' || character > '9' {
				t.Errorf("password %q contains non-digit %q", line, character)
			}
		}
	}
}

func TestRunInvalidCount(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--count", "0"}, &stdout, &stderr)
	if code == 0 {
		t.Fatal("run() code = 0, want non-zero")
	}
	if !strings.Contains(stderr.String(), "количество") {
		t.Fatalf("stderr = %q, want validation message", stderr.String())
	}
}

func TestRunNoEnabledGroups(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"--lowercase=false",
		"--uppercase=false",
		"--digits=false",
		"--symbols=false",
	}, &stdout, &stderr)
	if code == 0 {
		t.Fatal("run() code = 0, want non-zero")
	}
}

func TestRunHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() code = %d, want 0", code)
	}
	if !strings.Contains(stderr.String(), "Использование") {
		t.Fatalf("stderr = %q, want help text", stderr.String())
	}
}
