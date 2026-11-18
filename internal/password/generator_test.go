package password

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestGenerateLengthAndEnabledGroups(t *testing.T) {
	generator := New()
	opts := Options{
		Length:    32,
		Lowercase: true,
		Uppercase: true,
		Digits:    true,
		Symbols:   true,
	}

	for i := 0; i < 100; i++ {
		generated, err := generator.Generate(opts)
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}
		if len(generated) != opts.Length {
			t.Fatalf("len(Generate()) = %d, want %d", len(generated), opts.Length)
		}
		assertContainsAny(t, generated, Lowercase, "lowercase")
		assertContainsAny(t, generated, Uppercase, "uppercase")
		assertContainsAny(t, generated, Digits, "digit")
		assertContainsAny(t, generated, Symbols, "symbol")
	}
}

func TestGenerateUsesOnlyEnabledGroups(t *testing.T) {
	generator := New()
	generated, err := generator.Generate(Options{
		Length: 24,
		Digits: true,
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	for _, character := range generated {
		if !strings.ContainsRune(Digits, character) {
			t.Fatalf("Generate() produced disabled character %q", character)
		}
	}
}

func TestGenerateValidation(t *testing.T) {
	tests := []struct {
		name    string
		opts    Options
		wantErr error
	}{
		{
			name: "zero length",
			opts: Options{
				Lowercase: true,
			},
			wantErr: ErrInvalidLength,
		},
		{
			name: "no groups",
			opts: Options{
				Length: 16,
			},
			wantErr: ErrEmptyAlphabet,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New().Generate(test.opts)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Generate() error = %v, want %v", err, test.wantErr)
			}
		})
	}
}

func TestGenerateRejectsLengthBelowGroupCount(t *testing.T) {
	_, err := New().Generate(Options{
		Length:    2,
		Lowercase: true,
		Uppercase: true,
		Digits:    true,
	})
	if err == nil {
		t.Fatal("Generate() error = nil, want an error")
	}
}

func TestGeneratePropagatesRandomSourceError(t *testing.T) {
	_, err := NewWithReader(errorReader{}).Generate(Options{
		Length:    8,
		Lowercase: true,
	})
	if err == nil {
		t.Fatal("Generate() error = nil, want an error")
	}
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func assertContainsAny(t *testing.T, value, alphabet, group string) {
	t.Helper()
	if !strings.ContainsAny(value, alphabet) {
		t.Errorf("password %q does not contain a %s character", value, group)
	}
}
