package librsync_test

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/kelbyers/librsync-go"
)

func deltaSigToNew(sigName, baseName string) (out *bytes.Buffer, err error) {
	out = &bytes.Buffer{}
	sigFile, err := os.Open(sigName)
	defer sigFile.Close()
	if err != nil {
		return
	}
	n, err := os.Open(baseName)
	defer n.Close()
	if err != nil {
		return
	}

	output := bufio.NewWriter(out)

	sig, err := librsync.LoadSignature(sigFile)
	if err != nil {
		return
	}
	err = librsync.Delta(sig, n, output)

	output.Flush()

	return
}

func TestDelta(t *testing.T) {
	tests := []struct {
		name     string
		hashType string
		fileName string
		wantErr  bool
	}{
		{
			name:     "md4 01",
			hashType: "md4",
			fileName: "01",
		},
	}
	for _, tt := range tests {
		testDir := filepath.Join("testdir", "delta.input")
		sigDir := filepath.Join("testdir", "signature.input")
		sigName := filepath.Join(sigDir, tt.hashType, tt.fileName+".sig")
		expName := filepath.Join(testDir, tt.hashType, tt.fileName+".dlt")
		baseName := filepath.Join(testDir, tt.fileName+".in")
		t.Run(tt.name, func(t *testing.T) {
			delta, err := deltaSigToNew(sigName, baseName)
			if err != nil {
				t.Errorf("Error in deltaSigToNew: %v", err)
			}
			compareBuf2File(t, delta, expName)
		})
	}
}
