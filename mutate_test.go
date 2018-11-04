package librsync_test

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"testing"
)

func randInt(max int) int {
	r := rand.Int()
	return r % max
}

func xOffset(i, l, max int) int {
	o := i + l
	if o > max {
		return max
	}
	return o
}

func mutate(oldBytes []byte, seed int64, mutations int) []byte {
	rand.Seed(seed)
	nMuts := randInt(mutations) + 1

	for ; nMuts > 0; nMuts-- {
		inLength := len(oldBytes)
		fromOff := randInt(inLength)
		fromEnd := xOffset(fromOff, int(rand.Float32()*float32(randInt(inLength))), inLength)
		toOff := randInt(inLength)
		toEnd := xOffset(toOff, int(rand.Float32()*float32(randInt(inLength))), inLength)

		var chunk, pre, post []byte
		op := randInt(3)
		switch op {
		case 0:
			// copy and overwrite
			chunk = oldBytes[fromOff:fromEnd]
			pre = oldBytes[:toOff]
			post = oldBytes[toEnd:]
		case 1:
			// copy and insert
			chunk = oldBytes[fromOff:fromEnd]
			pre = oldBytes[:toOff]
			post = oldBytes[toOff:]
		case 2:
			// delete
			chunk = []byte{}
			pre = oldBytes[:toOff]
			post = oldBytes[toEnd:]
		}
		oldBytes = append([]byte{}, pre...)
		oldBytes = append(oldBytes, chunk...)
		oldBytes = append(oldBytes, post...)
	}
	return oldBytes
}

func TestMutate(t *testing.T) {
	o, _ := ioutil.ReadFile("mutate_test.go")
	for i := int64(0); i < 100; i++ {
		t.Logf("i: %d\n", i)
		n := mutate(o, i, 5)
		t.Logf("Changed: %v", bytes.Compare(o, n))
		if bytes.Compare(o, n) == 0 {
			t.Errorf("no change for i = %d\n", i)
		}
	}
}
