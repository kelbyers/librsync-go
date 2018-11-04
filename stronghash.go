package librsync

import (
	"fmt"
)

// block2hash converts a []byte into a value suitable for keying a map.
func block2hash(s []byte) string {
	return fmt.Sprintf("%s", s)
}

// StrongSignatureHashMap is used to map from a strong checksum to a position in
// the file.
type StrongSignatureHashMap struct {
	Strong map[string]int
}

// Get returns the position of the data block with strong checksum `k`
func (h *StrongSignatureHashMap) Get(k []byte) (int, bool) {
	key := block2hash(k)
	v, ok := h.Strong[key]
	return v, ok
}

// Set stores the position `l` for the block with strong checksum `k`
func (h *StrongSignatureHashMap) Set(k string, l int) {
	h.Strong[k] = l
}

// newStrongMap properly initializes a new StrongSignatureHashMap so that it can
// be updated, and returns the new StrongSignatureHashMap
func newStrongMap() StrongSignatureHashMap {
	st := make(map[string]int)
	h := StrongSignatureHashMap{Strong: st}
	return h
}
