package librsync_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	librsync "github.com/kelbyers/librsync-go"
)

func compareSig(t *testing.T, b *bytes.Buffer, e string) {
	expected, _ := ioutil.ReadFile(e)
	if !bytes.Equal(b.Bytes(), expected) {
		t.Errorf("Signatures don't match for %q", e)
	}
}

func TestSignature(t *testing.T) {
	tests := []struct {
		hashType string
		sigType  librsync.MagicNumber
		inFile   string
	}{
		{hashType: "md4", sigType: librsync.MD4_SIG_MAGIC,
			inFile: "01"},
	}
	for _, tt := range tests {
		name := tt.hashType + ":" + tt.inFile
		inFile := filepath.Join("tests", "signature.input", tt.inFile+".in")
		expectFile := filepath.Join("tests", "signature.input", tt.hashType, tt.inFile+".sig")
		var blockSize = uint32(2048)
		var sumSize uint32
		switch tt.hashType {
		case "md4":
			sumSize = uint32(16)
		case "blake2":
			sumSize = uint32(32)
		}
		t.Run(name, func(t *testing.T) {
			fmt.Printf("Input: %q\n", inFile)
			fmt.Printf("Comparison: %q\n", expectFile)

			basis, err := os.Open(inFile)
			defer basis.Close()
			if err != nil {
				t.Errorf("problem opening %q: %v", inFile, err)
			}

			var sigBuf bytes.Buffer
			sigWriter := bufio.NewWriter(&sigBuf)

			_, _ = librsync.Signature(basis, sigWriter, blockSize, sumSize, tt.sigType)
			compareSig(t, &sigBuf, expectFile)
		})
	}
}
