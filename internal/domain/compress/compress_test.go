package compress

import (
	"testing"
)

func TestCompressString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "empty string",
			input:   "",
			wantErr: false,
		},
		{
			name:    "short string",
			input:   "test",
			wantErr: false,
		},
		{
			name:    "longer string",
			input:   "this is a longer string that should compress reasonably well when repeated multiple times",
			wantErr: false,
		},
		{
			name:    "repeating content",
			input:   "abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := CompressString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompressString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Make sure we got a result if no error
			if err == nil && compressed == "" {
				t.Errorf("CompressString() returned empty result for input %q", tt.input)
			}

			// Verify we can decompress back to original
			decompressed, err := DecompressString(compressed)
			if err != nil {
				t.Errorf("Failed to decompress: %v", err)
			}

			if decompressed != tt.input {
				t.Errorf("Round-trip failed: got %q, want %q", decompressed, tt.input)
			}
		})
	}
}
