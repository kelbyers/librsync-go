package librsync

// SignatureHashMap maps a weak checksum to a hash of strong checksums for
// blocks that have the same weak checksum. In order to fully match a block of
// data, the weak checksum is used to get a StrongSignatureHashMap, and then the
// strong checksum is used to get the position of the block that matches both
// the weak sum and the strong sum.
type SignatureHashMap struct {
	Weak map[uint32]StrongSignatureHashMap
}

// Get returns the StrongSignatureHashMap for the provided weak checksum. In
// addition to the StrongSignatureHashMap, a boolean value will be returned,
// with `true` indicating that the weak checksum was found, and `false`
// indicating that it was not found.
func (h *SignatureHashMap) Get(k uint32) (StrongSignatureHashMap, bool) {
	v, ok := h.Weak[k]
	return v, ok
}

// Set associates weak sum `k` with the StrongSignatureHashMap `v`.
func (h *SignatureHashMap) Set(k uint32, v StrongSignatureHashMap) {
	h.Weak[k] = v
}

// UpdateBlock adds the signature mappings for the data block at position
// `l`, having weak checkum `w` and strong checksum `s`.
func (h *SignatureHashMap) UpdateBlock(w uint32, s []byte, l int) {
	st, ok := h.Get(w)
	if !ok {
		st = newStrongMap()
		h.Set(w, st)
	}
	st.Set(s, l)
}

// NewSignatureHashMap initializes a new SignatureHashMap so that
// StrongSignatureHashMaps can be added to it, and returns the new
// SignatureHashMap
func NewSignatureHashMap() SignatureHashMap {
	w := make(map[uint32]StrongSignatureHashMap)
	return SignatureHashMap{Weak: w}
}
