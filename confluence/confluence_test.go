package confluence

import (
	"testing"
)

func TestFormatBase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "valid file with extension",
			input:    "example-file_name.txt",
			expected: "example-file_name",
			wantErr:  false,
		},
		{
			name:     "file without extension",
			input:    "example-file_name",
			expected: "example-file_name",
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatBase(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatBase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("FormatBase() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
