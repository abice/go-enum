//go:build example
// +build example

package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuggyConstants(t *testing.T) {
	// Test that constants have correct values
	assert.Equal(t, Buggy(0), BuggyA, "BuggyA should have value 0")
	assert.Equal(t, Buggy(2), BuggyB, "BuggyB should have value 2")
	assert.Equal(t, Buggy(1), BuggyC, "BuggyC should have value 1")
}

func TestBuggyString(t *testing.T) {
	tests := []struct {
		name     string
		value    Buggy
		expected string
	}{
		{
			name:     "BuggyA",
			value:    BuggyA,
			expected: "A",
		},
		{
			name:     "BuggyB",
			value:    BuggyB,
			expected: "B",
		},
		{
			name:     "BuggyC",
			value:    BuggyC,
			expected: "C",
		},
		{
			name:     "Invalid value",
			value:    Buggy(99),
			expected: "Buggy(99)",
		},
		{
			name:     "Invalid value 3",
			value:    Buggy(3),
			expected: "Buggy(3)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.value.String())
		})
	}
}

func TestBuggyIsValid(t *testing.T) {
	tests := []struct {
		name     string
		value    Buggy
		expected bool
	}{
		{
			name:     "BuggyA is valid",
			value:    BuggyA,
			expected: true,
		},
		{
			name:     "BuggyB is valid",
			value:    BuggyB,
			expected: true,
		},
		{
			name:     "BuggyC is valid",
			value:    BuggyC,
			expected: true,
		},
		{
			name:     "Invalid value 99",
			value:    Buggy(99),
			expected: false,
		},
		{
			name:     "Invalid value 3",
			value:    Buggy(3),
			expected: false,
		},
		{
			name:     "Invalid value 3",
			value:    Buggy(3),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.value.IsValid())
		})
	}
}

func TestParseBuggy(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      Buggy
		errorExpected bool
	}{
		{
			name:          "Parse A",
			input:         "A",
			expected:      BuggyA,
			errorExpected: false,
		},
		{
			name:          "Parse B",
			input:         "B",
			expected:      BuggyB,
			errorExpected: false,
		},
		{
			name:          "Parse C",
			input:         "C",
			expected:      BuggyC,
			errorExpected: false,
		},
		{
			name:          "Parse invalid string",
			input:         "D",
			expected:      Buggy(0),
			errorExpected: true,
		},
		{
			name:          "Parse empty string",
			input:         "",
			expected:      Buggy(0),
			errorExpected: true,
		},
		{
			name:          "Parse lowercase a",
			input:         "a",
			expected:      Buggy(0),
			errorExpected: true,
		},
		{
			name:          "Parse numeric string",
			input:         "1",
			expected:      Buggy(0),
			errorExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseBuggy(tt.input)

			if tt.errorExpected {
				require.Error(t, err, "Expected error for input: %s", tt.input)
				assert.Contains(t, err.Error(), "not a valid Buggy", "Error should contain 'not a valid Buggy'")
				assert.Equal(t, tt.expected, result, "Result should be zero value on error")
			} else {
				require.NoError(t, err, "Unexpected error for input: %s", tt.input)
				assert.Equal(t, tt.expected, result, "Parsed value should match expected")
			}
		})
	}
}

func TestBuggyErrorType(t *testing.T) {
	// Test that ErrInvalidBuggy is properly defined
	assert.NotNil(t, ErrInvalidBuggy)
	assert.Equal(t, "not a valid Buggy", ErrInvalidBuggy.Error())

	// Test that ParseBuggy returns an error that wraps ErrInvalidBuggy
	_, err := ParseBuggy("invalid")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidBuggy)
}

func TestBuggyEdgeCases(t *testing.T) {
	t.Run("Zero value behavior", func(t *testing.T) {
		var b Buggy
		// Zero value should be 0, which is BuggyA
		assert.Equal(t, Buggy(0), b)
		assert.Equal(t, BuggyA, b)
		assert.True(t, b.IsValid())
		assert.Equal(t, "A", b.String())
	})

	t.Run("Non-sequential values", func(t *testing.T) {
		// Test that the enum handles non-sequential values correctly
		// A=0, B=2, C=1 - so value 1 should map to C, not B
		assert.Equal(t, "C", Buggy(1).String())
		assert.Equal(t, "B", Buggy(2).String())
		assert.True(t, Buggy(1).IsValid())
		assert.True(t, Buggy(2).IsValid())
	})
}

func BenchmarkBuggyString(b *testing.B) {
	values := []Buggy{BuggyA, BuggyB, BuggyC, Buggy(99)}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			_ = v.String()
		}
	}
}

func BenchmarkParseBuggy(b *testing.B) {
	inputs := []string{"A", "B", "C", "invalid"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, input := range inputs {
			_, _ = ParseBuggy(input)
		}
	}
}

func BenchmarkBuggyIsValid(b *testing.B) {
	values := []Buggy{BuggyA, BuggyB, BuggyC, Buggy(99)}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			_ = v.IsValid()
		}
	}
}
