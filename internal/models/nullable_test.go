package models

import (
	"testing"
)

func TestNullStringMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    NullString
		expected string
	}{
		{
			name:     "Valid true",
			input:    NullString{String: "test", Valid: true},
			expected: `"test"`,
		}, {
			name:     "Valid false",
			input:    NullString{String: "", Valid: false},
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.MarshalJSON()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(got))
			}
		})
	}
}
