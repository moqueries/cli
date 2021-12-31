package hash

import (
	"encoding/binary"

	"github.com/davegardnerisme/deephash"
)

// Hash stores the hash of a source object
type Hash uint64

// DeepHash walks the src parameter and produces a hash
func DeepHash(src interface{}) Hash {
	b := deephash.Hash(src)
	if len(b) < 8 {
		newB := make([]byte, 8)
		copy(newB, b)
		b = newB
	}
	return Hash(binary.LittleEndian.Uint64(b))
}
