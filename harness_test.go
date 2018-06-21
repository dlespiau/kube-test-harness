package harness

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveDirectory(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)

	tests := []struct {
		in       string
		expected string
	}{
		{"", wd},
		{"/foo/bar", "/foo/bar"},
		{"foo", filepath.Join(wd, "foo")},
		{"foo/bar", filepath.Join(wd, "foo/bar")},
	}

	for _, test := range tests {
		dir, err := resolveDirectory(test.in)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, dir)
	}
}
