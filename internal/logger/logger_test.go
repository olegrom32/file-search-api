package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		prefix     string
		wantPrefix string
		wantErr    string
	}{
		{
			name:       "given info prefix, then should return the logger",
			prefix:     "info",
			wantPrefix: "INFO: ",
		},
		{
			name:       "given debug prefix, then should return the logger",
			prefix:     "debug",
			wantPrefix: "DEBUG: ",
		},
		{
			name:       "given error prefix, then should return the logger",
			prefix:     "error",
			wantPrefix: "ERROR: ",
		},
		{
			name:    "given unsupported prefix, then should return error",
			prefix:  "alert",
			wantErr: "invalid logger prefix:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.prefix)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			}

			if tt.wantPrefix != "" {
				assert.Equal(t, tt.wantPrefix, got.Prefix())
			}
		})
	}
}
