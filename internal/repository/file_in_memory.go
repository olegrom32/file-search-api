package repository

import (
	"math"
	"sort"

	"github.com/olegrom32/file-search-api/internal"
)

// FileInMemory is an in-memory implementation
type FileInMemory struct {
	file   []int
	margin float64
}

func NewFileInMemory(file []int, margin float64) *FileInMemory {
	return &FileInMemory{
		file:   file,
		margin: margin,
	}
}

// FindByValue searches and returns a index of the given value.
// The implementation uses sort.Search which is binary search function.
// Binary search is the most optimal algo for searching in sorted lists.
// I thought there is no point in implementing the binary search myself as then I will
// have to pretty much just copy the contents of sort.Search function here.
func (r *FileInMemory) FindByValue(value int) (int, error) {
	potentialIdx := -1
	bestDiff := math.MaxInt

	// We are looking not only for exact value (which is the best result), we are also looking for
	// any value within the configured margin.
	margin := int(math.Round(float64(value) * r.margin))

	sort.Search(len(r.file), func(i int) bool {
		// Check if current value is withing the acceptable range.
		if r.file[i] <= value+margin && r.file[i] >= value-margin {
			diff := r.file[i] - value
			if diff < 0 {
				diff *= -1
			}

			// Try to return the value as close as possible to the target value.
			// If the value is equal to the target value, the diff will be 0 making it the best option to return.
			if diff < bestDiff {
				potentialIdx = i
				bestDiff = diff
			}
		}

		return r.file[i] >= value
	})

	// Return the best option
	if potentialIdx > -1 {
		return potentialIdx, nil
	}

	// Nothing found
	return 0, internal.ErrNotFound
}
