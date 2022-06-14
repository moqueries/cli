// Package hash implements a deep hash mechanism so that arbitrary types can be
// used as a key in hash maps
package hash

import (
	"encoding/binary"

	"github.com/davegardnerisme/deephash"
)

const hashBytes = 8

// Hash stores the hash of a source object
type Hash uint64

// DeepHash walks the src parameter and produces a hash
func DeepHash(src interface{}) Hash {
	b := deephash.Hash(src)
	if len(b) < hashBytes {
		newB := make([]byte, hashBytes)
		copy(newB, b)
		b = newB
	}
	return Hash(binary.LittleEndian.Uint64(b))
}
