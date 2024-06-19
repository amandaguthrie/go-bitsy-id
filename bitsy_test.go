package bitsy

import (
	"fmt"
	"strconv"
	"testing"
)

var idSizes = []struct {
	size int
}{
	{size: 4},
	{size: 8},
	{size: 12},
	{size: 16},
	{size: 20},
	{size: 24},
	{size: 28},
	{size: 32},
	{size: 64},
}

var multipliers = []struct {
	multiplier int
}{
	{multiplier: defaultPoolSizeMultiplier},
}

func TestLoadPool(t *testing.T) {
	t.Run("offset + Size > pool length", func(t *testing.T) {
		t.Run("poolLen < Size", func(t *testing.T) {
			// Set up pool Size to allow loadPool to fall into the poolLen < Size condition
			poolLen := len(pool)
			newPoolSize := poolLen - offset
			createPool(newPoolSize)

			// Test
			newPoolSize = poolLen + 12
			loadPool(newPoolSize)

			// The offset should match the new pool Size since createPool is called in loadPool, which resets offset to 0
			if offset != newPoolSize {
				t.Fatal("Offset not adjusted")
			}
		})
		t.Run("poolLen > Size", func(t *testing.T) {
			// Set up pool Size to allow loadPool to fall into the poolLen > Size condition
			poolLen := len(pool)
			newPoolSize := poolLen - offset + 1
			createPool(newPoolSize)

			// Test
			poolLen = len(pool)
			offsetBefore := offset
			newPoolSize = poolLen - 12
			loadPool(newPoolSize)

			// The offset should match offset before + Size from loadPool
			if offset != offsetBefore+newPoolSize {
				t.Fatal("Offset not adjusted")
			}
		})
	})
}

func TestBitsyID(t *testing.T) {
	size := 21

	id := ID(&IDParams{Size: size})
	id2 := ID(&IDParams{Size: size})
	ID(&IDParams{Size: size})

	t.Run("should generate a ID with length"+strconv.Itoa(size), func(t *testing.T) {
		if len(id) != size || len(id2) != size {
			t.Fatalf("BitsyId did not generate an ID with length %d.", size)
		}
	})
	t.Run("should generate different BitsyIDs when the function is run twice", func(t *testing.T) {
		if id == id2 {
			t.Fatal("ID generated the same ID twice in a row.")
		}
	})

	t.Run("should use custom Multiplier if provided and valid", func(t *testing.T) {
		customMultiplier := 256
		ID(&IDParams{Size: size, Multiplier: customMultiplier})
		if multiplier != customMultiplier {
			t.Fatal("Custom Alphabet not used.")
		}
	})

	t.Run("should use custom Alphabet if provided", func(t *testing.T) {
		customAlpha := []rune("a")
		sz := 3
		alphaID := ID(&IDParams{Size: sz, Alphabet: customAlpha})
		for i := 0; i < sz; i++ {
			if string(alphaID[i]) != "a" {
				t.Fatal("Custom Alphabet not used.")
			}
		}
	})

}

func BenchmarkBitsyID(b *testing.B) {
	for _, m := range multipliers {
		for _, v := range idSizes {
			size := v.size
			b.Run(fmt.Sprintf("%d-%d", m.multiplier, size), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					ID(&IDParams{Size: size, Multiplier: m.multiplier})
				}
			})
		}
	}
}
