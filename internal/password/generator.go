// Package password implements a cryptographically secure password generator.
package password

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"
)

const (
	Lowercase = "abcdefghijklmnopqrstuvwxyz"
	Uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits    = "0123456789"
	Symbols   = "!@#$%^&*()-_=+[]{}:,.?"
)

var (
	ErrInvalidLength = errors.New("password length must be greater than zero")
	ErrEmptyAlphabet = errors.New("at least one character set must be enabled")
)

// Options controls password generation.
type Options struct {
	Length    int
	Lowercase bool
	Uppercase bool
	Digits    bool
	Symbols   bool
}

// Generator generates passwords using the supplied source of random bytes.
type Generator struct {
	random io.Reader
}

// New returns a generator backed by crypto/rand.Reader.
func New() *Generator {
	return &Generator{random: rand.Reader}
}

// NewWithReader returns a generator backed by random.
// It is primarily useful for deterministic tests.
func NewWithReader(random io.Reader) *Generator {
	return &Generator{random: random}
}

// Generate creates one password matching opts.
// Every enabled character group is guaranteed to occur at least once.
func (g *Generator) Generate(opts Options) (string, error) {
	groups := enabledGroups(opts)
	if opts.Length <= 0 {
		return "", ErrInvalidLength
	}
	if len(groups) == 0 {
		return "", ErrEmptyAlphabet
	}
	if opts.Length < len(groups) {
		return "", fmt.Errorf("password length (%d) must be at least the number of enabled character sets (%d)", opts.Length, len(groups))
	}
	if g == nil || g.random == nil {
		return "", errors.New("random source is not configured")
	}

	alphabet := strings.Join(groups, "")
	result := make([]byte, opts.Length)

	// Add one random character from each enabled group first.
	for i, group := range groups {
		character, err := g.randomByte(group)
		if err != nil {
			return "", fmt.Errorf("read random data: %w", err)
		}
		result[i] = character
	}

	// Fill all remaining positions from the combined alphabet.
	for i := len(groups); i < len(result); i++ {
		character, err := g.randomByte(alphabet)
		if err != nil {
			return "", fmt.Errorf("read random data: %w", err)
		}
		result[i] = character
	}

	// Shuffle so the required characters are not always at the beginning.
	for i := len(result) - 1; i > 0; i-- {
		j, err := g.randomIndex(i + 1)
		if err != nil {
			return "", fmt.Errorf("read random data: %w", err)
		}
		result[i], result[j] = result[j], result[i]
	}

	return string(result), nil
}

func enabledGroups(opts Options) []string {
	groups := make([]string, 0, 4)
	if opts.Lowercase {
		groups = append(groups, Lowercase)
	}
	if opts.Uppercase {
		groups = append(groups, Uppercase)
	}
	if opts.Digits {
		groups = append(groups, Digits)
	}
	if opts.Symbols {
		groups = append(groups, Symbols)
	}
	return groups
}

func (g *Generator) randomByte(alphabet string) (byte, error) {
	index, err := g.randomIndex(len(alphabet))
	if err != nil {
		return 0, err
	}
	return alphabet[index], nil
}

func (g *Generator) randomIndex(limit int) (int, error) {
	value, err := rand.Int(g.random, big.NewInt(int64(limit)))
	if err != nil {
		return 0, err
	}
	return int(value.Int64()), nil
}
