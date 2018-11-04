package librsync

import "github.com/dchest/siphash"

const (
	// generated by splitting the md5 sum of "hashmap"
	sipHashKey1 = 0xdda7806a4847ec61
	sipHashKey2 = 0xb5940c2623a5aabd
)

// block2hash converts a []byte into a value suitable for keying a map. This
// code is based on the method used in `github.com/cornelk/hashmap` for making
// []byte into a hashable key
func block2hash(s []byte) uintptr {
	return uintptr(siphash.Hash(sipHashKey1, sipHashKey2, s))
}

// StrongSignatureHashMap is used to map from a strong checksum to a position in
// the file.
type StrongSignatureHashMap struct {
	Strong map[uintptr]int
}

// Get returns the position of the data block with strong checksum `k`
func (h *StrongSignatureHashMap) Get(k []byte) (int, bool) {
	key := block2hash(k)
	v, ok := h.Strong[key]
	return v, ok
}

// Set stores the position `l` for the block with strong checksum `k`
func (h *StrongSignatureHashMap) Set(k uintptr, l int) {
	h.Strong[k] = l
}

// newStrongMap properly initializes a new StrongSignatureHashMap so that it can
// be updated, and returns the new StrongSignatureHashMap
func newStrongMap() StrongSignatureHashMap {
	st := make(map[uintptr]int)
	h := StrongSignatureHashMap{Strong: st}
	return h
}
