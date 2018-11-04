package librsync_test

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"testing"

	librsync "github.com/kelbyers/librsync-go"
)

func patchBaseFromDiff(baseName, deltaName string) (out *bytes.Buffer, err error) {
	out = &bytes.Buffer{}
	base, err := os.Open(baseName)
	defer base.Close()
	if err != nil {
		return
	}
	delta, err := os.Open(deltaName)
	defer delta.Close()
	if err != nil {
		return
	}
	w := bufio.NewWriter(out)
	librsync.Patch(base, delta, w)
	w.Flush()

	return
}

func TestPatchFromNull(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"01"},
		{"02"},
		{"03"},
	}
	for _, tt := range tests {
		testDir := filepath.Join("testdir", "patch.input")
		deltaName := filepath.Join(testDir, tt.name+".delta")
		expectFile := filepath.Join(testDir, tt.name+".expect")
		t.Run(tt.name, func(t *testing.T) {
			out, err := patchBaseFromDiff(os.DevNull, deltaName)
			if err != nil {
				t.Errorf("Problem opening files for patching: %s", err)
			}
			compareBuf2File(t, out, expectFile)
		})
	}
}
