package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/olegrom32/file-search-api/internal"
)

func TestFileInMemory_FindByValue(t *testing.T) {
	r := NewFileInMemory([]int{
		100, 200, 300, 500, 502, 550, 1000, 2000,
	}, 0.1)

	tests := []struct {
		name    string
		value   int
		want    int
		wantErr error
	}{
		{
			name:    "given an input of 500, then should return the index of 500",
			value:   500,
			want:    3,
			wantErr: nil,
		},
		{
			name:    "given an input of 501, then should return the index of 550",
			value:   501,
			want:    4,
			wantErr: nil,
		},
		{
			name:    "given an input of 551, then should return the index of 550",
			value:   551,
			want:    5,
			wantErr: nil,
		},
		{
			name:    "given an input of 1999, then should return the index of 2000",
			value:   1999,
			want:    7,
			wantErr: nil,
		},
		{
			name:    "given an input of 2001, then should return the index of 2000",
			value:   2001,
			want:    7,
			wantErr: nil,
		},
		{
			name:    "given an input of 700, then should return error",
			value:   700,
			want:    0,
			wantErr: internal.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.FindByValue(tt.value)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
