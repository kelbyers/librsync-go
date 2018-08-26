package librsync

// Test cases based on tests in the real librsync source
// https://github.com/librsync/librsync/blob/47ed48c62dfdc1a945db9e8e0b49a768dd376ea0/tests/rollsum_test.c

import (
	"fmt"
	"testing"
)

// testAssert makes it slightly less verbose to compare an expected value with a
// test result and then fail with a nice message
func testAssert(t *testing.T, expected, got interface{}, msgAndArgs ...interface{}) {
	if expected != got {
		var msg string
		if len(msgAndArgs) > 0 {
			format := msgAndArgs[0].(string)
			args := msgAndArgs[1:]
			msg = ": " + fmt.Sprintf(format, args...)
		}
		t.Errorf("Expected %v, got %v%s", expected, got, msg)
	}
}

func TestWeakChecksum(t *testing.T) {
	tests := []struct {
		name  string
		count int
		want  uint32
	}{
		{name: "first 256 bytes", count: 256, want: 0x3a009e80},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := make([]byte, tt.count)
			for i := 0; i < tt.count; i++ {
				p[i] = byte(i)
			}
			got := WeakChecksum(p)
			testAssert(t, tt.want, got)
		})
	}
}

func TestNewRollsum(t *testing.T) {
	tests := []struct {
		name   string
		count  uint64
		s1, s2 uint16
		digest uint32
	}{
		{name: "NewRollsum", count: 0, s1: 0, s2: 0, digest: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRollsum()
			testAssert(t, tt.s1, got.s1, "s1")
			testAssert(t, tt.s2, got.s2, "s2")
			testAssert(t, tt.count, got.count, "count")
			testAssert(t, tt.digest, got.Digest(), "digest")
		})
	}
}

func TestRollsum_Update(t *testing.T) {
	r := Rollsum{}
	tests := []struct {
		name   string
		count  int
		digest uint32
	}{
		{name: "first 256 bytes", count: 256, digest: 0x3a009e80},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := make([]byte, tt.count)
			for i := 0; i < tt.count; i++ {
				p[i] = byte(i)
			}
			r.Update(p)
			testAssert(t, tt.digest, r.Digest())
		})
	}
}

func TestRollsum_Rotate(t *testing.T) {
	r := Rollsum{4, 130, 320}
	type args struct {
		out byte
		in  byte
	}
	tests := []struct {
		name   string
		args   []args
		count  uint64
		digest uint32
	}{
		{name: "[1,2,3,4]", args: []args{{0, 4}}, count: 4, digest: 0x014a0086},
		{name: "[4,5,6,7]", args: []args{{1, 5}, {2, 6}, {3, 7}}, count: 4, digest: 0x01680092},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.args {
				r.Rotate(v.out, v.in)
			}
			testAssert(t, tt.count, r.count)
			testAssert(t, tt.digest, r.Digest())
		})
	}
}

func TestRollsum_Rollin(t *testing.T) {
	r := Rollsum{}
	tests := []struct {
		name   string
		in     []byte
		count  uint64
		digest uint32
	}{
		{name: "Empty 0", in: []byte{0}, count: 1, digest: uint32(0x001f001f)},
		{name: "[1,2,3]", in: []byte{1, 2, 3}, count: 4, digest: uint32(0x01400082)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, i := range tt.in {
				r.Rollin(i)
			}
			testAssert(t, tt.count, r.count)
			testAssert(t, tt.digest, r.Digest())
		})
	}
}

func TestRollsum_Rollout(t *testing.T) {
	r := Rollsum{4, 146, 360}
	tests := []struct {
		name   string
		out    []byte
		count  uint64
		digest uint32
	}{
		{name: "[5,6,7]", out: []byte{4}, count: 3, digest: 0x00dc006f},
		{name: "[]", out: []byte{5, 6, 7}, count: 0, digest: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, out := range tt.out {
				r.Rollout(out)
			}
			testAssert(t, tt.count, r.count, "count")
			testAssert(t, tt.digest, r.Digest(), "digest")
		})
	}
}

func TestRollsum_Digest(t *testing.T) {
	tests := []struct {
		name string
		r    Rollsum
		want uint32
	}{
		{name: "new Rollsum digest", r: NewRollsum(), want: uint32(0)},
		{name: "rollin+rotate", r: Rollsum{4, 146, 360}, want: 0x01680092},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.r.Digest()
			testAssert(t, tt.want, got)
		})
	}
}

func TestRollsum_Reset(t *testing.T) {
	tests := []struct {
		name   string
		r      Rollsum
		count  uint64
		s1, s2 uint16
	}{
		{name: "Already empty"},
		{name: "Not empty", r: Rollsum{4, 146, 360}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Reset()
			testAssert(t, tt.count, tt.r.count, "count")
			testAssert(t, tt.s1, tt.r.s1, "s1")
			testAssert(t, tt.s2, tt.r.s2, "s2")
		})
	}
}

func benchmarkRollsum_Update(count int, b *testing.B) {
	r := Rollsum{}
	p := makeByteSlice(count)
	for n := 0; n < b.N; n++ {
		r.Reset()
		r.Update(p)
	}
}

func makeByteSlice(count int) []byte {
	p := make([]byte, count)
	for i := 0; i < count; i++ {
		p[i] = byte(i)
	}
	return p
}

func BenchmarkRollsum_Update256(b *testing.B)  { benchmarkRollsum_Update(256, b) }
func BenchmarkRollsum_Update1024(b *testing.B) { benchmarkRollsum_Update(1024, b) }

var benchRollsum Rollsum

func BenchmarkRollsum_UpdateComplete(b *testing.B) {
	p := makeByteSlice(256)
	r := NewRollsum()
	for n := 0; n < b.N; n++ {
		r.Reset()
		r.Update(p)
	}
	benchRollsum = r
}

func benchmarkRollsum_Rollin(count int, b *testing.B) {
	r := NewRollsum()
	p := makeByteSlice(count)
	for n := 0; n < b.N; n++ {
		r.Reset()
		for _, in := range p {
			r.Rollin(in)
		}
	}
}

func BenchmarkRollsum_Rollin1(b *testing.B)   { benchmarkRollsum_Rollin(1, b) }
func BenchmarkRollsum_Rollin2(b *testing.B)   { benchmarkRollsum_Rollin(2, b) }
func BenchmarkRollsum_Rollin5(b *testing.B)   { benchmarkRollsum_Rollin(5, b) }
func BenchmarkRollsum_Rollin10(b *testing.B)  { benchmarkRollsum_Rollin(10, b) }
func BenchmarkRollsum_Rollin256(b *testing.B) { benchmarkRollsum_Rollin(256, b) }

func BenchmarkRollsum_RollinComplete(b *testing.B) {
	r := NewRollsum()
	p := makeByteSlice(10)
	for n := 0; n < b.N; n++ {
		r.Reset()
		for _, in := range p {
			r.Rollin(in)
		}
	}
	benchRollsum = r
}

func benchmarkRollsumRollout(count int, b *testing.B) {
	r := NewRollsum()
	r.count = uint64(count)
	p := makeByteSlice(count)
	for n := 0; n < b.N; n++ {
		r.Reset()
		for _, out := range p {
			r.Rollout(out)
		}
	}
}

func BenchmarkRollsum_Rollout1(b *testing.B)   { benchmarkRollsumRollout(1, b) }
func BenchmarkRollsum_Rollout2(b *testing.B)   { benchmarkRollsumRollout(2, b) }
func BenchmarkRollsum_Rollout5(b *testing.B)   { benchmarkRollsumRollout(5, b) }
func BenchmarkRollsum_Rollout10(b *testing.B)  { benchmarkRollsumRollout(10, b) }
func BenchmarkRollsum_Rollout256(b *testing.B) { benchmarkRollsumRollout(256, b) }

func BenchmarkRollsum_RolloutComplete(b *testing.B) {
	r := NewRollsum()
	count := 10
	r.count = uint64(count)
	p := makeByteSlice(count)
	for n := 0; n < b.N; n++ {
		r.Reset()
		for _, out := range p {
			r.Rollout(out)
		}
	}
	benchRollsum = r
}

var benchDigest uint32

func BenchmarkRollsum_DigestComplete(b *testing.B) {
	var d uint32
	rollsum := Rollsum{s1: uint16(1111), s2: uint16(2222)}
	for n := 0; n < b.N; n++ {
		d = rollsum.Digest()
	}
	benchDigest = d
}
