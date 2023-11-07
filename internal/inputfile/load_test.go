package inputfile

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/olegrom32/file-search-api/internal"
)

type testInputFile struct {
	io.ReadSeeker
}

func (f testInputFile) Close() error {
	return nil
}

func testFile(contents string) testInputFile {
	return testInputFile{ReadSeeker: strings.NewReader(contents)}
}

func TestLoad(t *testing.T) {
	t.Run("given valid input, then should process the file", func(t *testing.T) {
		got, err := Load(testFile("100\n200\n300\n"))
		require.NoError(t, err)

		assert.Equal(t, []int{100, 200, 300}, got)
	})

	t.Run("given valid and invalid lines in the input, then should return error", func(t *testing.T) {
		_, err := Load(testFile("100\n200\ninvalid\n"))
		assert.ErrorIs(t, err, internal.ErrInvalidInputFile)
	})

	t.Run("given an empty input, then should return error", func(t *testing.T) {
		_, err := Load(testFile(""))
		assert.ErrorIs(t, err, internal.ErrInvalidInputFile)
	})

	t.Run("given a file with no terminating newline, then should return error", func(t *testing.T) {
		// For the test task purposes, let's support only POSIX-compliant files (https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap03.html#tag_03_206)
		_, err := Load(testFile("123"))
		assert.ErrorIs(t, err, internal.ErrInvalidInputFile)
	})
}
