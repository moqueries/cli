package hash_test

import (
	"testing"

	"moqueries.org/cli/hash"
)

func TestNil(t *testing.T) {
	// ASSEMBLE

	// ACT
	//nolint:ifshort // short syntax here blurs ACT/ASSERT
	h := hash.DeepHash(nil)

	// ASSERT
	// The value doesn't really matter. Just make sure it doesn't panic.
	if h != 0 {
		t.Errorf("wanted 0, got %d", h)
	}
}
