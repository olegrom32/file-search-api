package inputfile

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/olegrom32/file-search-api/internal"
)

// Load loads the input file into the slice as requested in the task description
func Load(f io.ReadSeekCloser) ([]int, error) {
	defer func() {
		if err := f.Close(); err != nil {
			log.Print(err)
		}
	}()

	// In order to create a slice in an optimal way (i.e. allocating it once without the use of `append`),
	// let's first count the number of lines in the file. This requires reading the file twice,
	// which is hopefully still a good trade-off.
	n, err := lineCounter(f)
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, fmt.Errorf("file cannot be empty: %w", internal.ErrInvalidInputFile)
	}

	res := make([]int, n)

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to rewind the input file for processing: %w", err)
	}

	scanner := bufio.NewScanner(f)
	i := 0

	for scanner.Scan() {
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("input file contains invalid line at index %d: %s: %s: %w", i, scanner.Text(), err, internal.ErrInvalidInputFile)
		}

		res[i] = v
		i++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan input file: %w", err)
	}

	return res, nil
}

func lineCounter(r io.Reader) (int, error) {
	// Assuming 32kb buffer is enough to read a line.
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return count, nil
			}

			return 0, fmt.Errorf("failed to read from the input file: %w", err)
		}

		nNewlines := bytes.Count(buf[:c], lineSep)
		if nNewlines == 0 {
			return 0, fmt.Errorf("line does not end with a newline: %w", internal.ErrInvalidInputFile)
		}

		count += nNewlines
	}
}
