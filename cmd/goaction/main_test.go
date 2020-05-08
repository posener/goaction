package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPathRelDir(t *testing.T) {
	t.Parallel()

	tests := []struct{ path, want string }{
		{"", "./"},
		{"./", "./"},
		{"src/", "./src"},
		{"./src/", "./src"},
	}

	for _, tt := range tests {
		got, err := pathRelDir(tt.path)
		require.NoError(t, err)
		assert.Equal(t, tt.want, got)
	}
}
