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

func TestNullStringScan(t *testing.T) {
	tests := []struct {
		name           string
		input          interface{}
		expectedString string
		expectedValid  bool
	}{
		{
			name:           "Valid string",
			input:          "test",
			expectedString: "test",
			expectedValid:  true,
		},
		{
			name:           "Nil value",
			input:          nil,
			expectedString: "",
			expectedValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ns NullString
			err := ns.Scan(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			gotString := ns.String

			if gotString != tt.expectedString {
				t.Errorf("expected string %s, got %s", tt.expectedString, gotString)
			}

			gotValid := ns.Valid
			if gotValid != tt.expectedValid {
				t.Errorf("expected valid %v, got %v", tt.expectedValid, gotValid)
			}
		})
	}
}
