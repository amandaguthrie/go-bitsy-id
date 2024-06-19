// Package bitsy implements utility functions for generating Bitsy IDs.
package bitsy

import (
	"crypto/rand"
)

const (
	defaultSize               = 21
	defaultPoolSizeMultiplier = 128
)

var (
	// alphabet contains the runes used for generating the ID
	alphabet = []rune("useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict")
	// multiplier stores the size multiplier used for calculating pool length
	multiplier = defaultPoolSizeMultiplier
	// offset stores the last used index of random bytes in the pool used to generate an ID
	offset = 0
	// pool stores the random byte integers used for generating the ID
	pool []byte
)

// loadPool creates or refreshes the pool if required, and advances the offset.
func loadPool(size int) {
	poolLen := len(pool)

	if offset+size > poolLen {
		if poolLen < size {
			createPool(size)
		} else {
			populatePool()
		}
	}

	offset += size
}

// createPool makes and populates a new pool with the specified size and multiplier.
func createPool(size int) {
	pool = make([]byte, size*multiplier)
	populatePool()
}

// populatePool reads random bytes to fill the pool and resets the offset.
func populatePool() {
	_, _ = rand.Read(pool)
	resetOffset()
}

// randomSlice loads the pool and returns a slice calculated with the offset and size.
func randomSlice(size int) []byte {
	loadPool(size)
	return pool[offset-size : offset]
}

// resetOffset resets the offset to zero.
func resetOffset() {
	offset = 0
}

// IDParams define the parameter structure for the ID function.
//
// Alphabet
//   - Description: The array of characters mapped from runes to generate the final ID.
//   - Required: false
//   - Default: []rune("useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict")
//
// Multiplier
//   - Description: Multiplied with the size to determine the pool of bytes used to generate an ID. The pool has
//     enough characters to generate multiple IDs before it regenerates. Benchmark testing showed that
//     generally generating larger pools rather than smaller pools for each ID was more efficient.
//   - Required: false
//   - Minimum: 1
//   - Default: 128
//
// Size
//   - Description: The number of characters in a generated ID.
//   - Required: false
//   - Minimum: 1
//   - Default: 21
type IDParams struct {
	Alphabet   []rune
	Multiplier int
	Size       int
}

// ID generates a random string of the specified size.
// Use the arguments to customize the alphabet and size of the generated IDs.
func ID(p *IDParams) string {
	var sz = defaultSize

	if len(p.Alphabet) > 0 {
		alphabet = p.Alphabet
	}
	alphaLen := len(alphabet)
	if p.Size >= 1 {
		sz = p.Size
	}
	if p.Multiplier >= 1 {
		multiplier = p.Multiplier
	}
	thisPool := randomSlice(sz)
	var id = make([]rune, sz)

	for i := range id {
		id[i] = alphabet[thisPool[i]&byte(alphaLen-1)]
	}
	return string(id)
}
