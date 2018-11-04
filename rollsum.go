package librsync

// Rollsum holds the current state for calculating the rolling weak checksums of
// blocks
type Rollsum struct {
	count  uint64
	s1, s2 uint16
}

// RollsumCharOffset is a constant selected for performance reasons in the rolling
// checksum algorithm
const RollsumCharOffset = 31

// WeakChecksum returns the weak checksum for a block of data
func WeakChecksum(data []byte) uint32 {
	var sum Rollsum
	sum.Update(data)
	return sum.Digest()
}

// NewRollsum returns a new initialized Rollsum
func NewRollsum() Rollsum {
	return Rollsum{}
}

// Update the state for `r` with data from `p`
func (r *Rollsum) Update(p []byte) {
	l := len(p)

	for n := 0; n < l; {
		if n+15 < l {
			for i := 0; i < 16; i++ {
				r.s1 += uint16(p[n+i])
				r.s2 += r.s1
			}
			n += 16
		} else {
			r.s1 += uint16(p[n])
			r.s2 += r.s1
			n++
		}
	}

	r.s1 += uint16(l * RollsumCharOffset)
	r.s2 += uint16(((l * (l + 1)) / 2) * RollsumCharOffset)
	r.count += uint64(l)
}

// Rotate replaces old byte `out` with new byte `in`
func (r *Rollsum) Rotate(out, in byte) {
	r.s1 += uint16(in - out)
	r.s2 += r.s1 - uint16(r.count)*(uint16(out)+uint16(RollsumCharOffset))
}

// Rollin adds `in` to the rolling checksum
func (r *Rollsum) Rollin(in byte) {
	r.s1 += uint16(in) + uint16(RollsumCharOffset)
	r.s2 += r.s1
	r.count++
}

// Rollout removes `out` from the rolling checksum
func (r *Rollsum) Rollout(out byte) {
	r.s1 -= uint16(out) + uint16(RollsumCharOffset)
	r.s2 -= uint16(r.count) * (uint16(out) + uint16(RollsumCharOffset))
	r.count--
}

// Digest combines the rolling checksum components and returns the actual
// checksum.
func (r *Rollsum) Digest() uint32 {
	return (uint32(r.s2) << 16) | (uint32(r.s1) & 0xffff)
}

// Reset re-zeros the rolling checksum to the zero value
func (r *Rollsum) Reset() {
	r.count = 0
	r.s1 = 0
	r.s2 = 0
}
