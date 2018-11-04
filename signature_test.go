package librsync_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	librsync "github.com/kelbyers/librsync-go"
)

func compareBuf2File(t *testing.T, b *bytes.Buffer, e string) {
	expected, _ := ioutil.ReadFile(e)
	got := b.Bytes()
	if !bytes.Equal(got, expected) {
		t.Errorf("Signatures don't match for %q", e)
	}
}

func sigFromFile(
	name string, sigType librsync.MagicNumber, sumSize, blockSize uint32,
) (*librsync.SignatureType, *bytes.Buffer, error) {
	basis, err := os.Open(name)
	defer basis.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to open basis file(%s) - %v",
			name, err)
	}
	sigBuf := &bytes.Buffer{}
	sigWriter := bufio.NewWriter(sigBuf)
	signature, err := librsync.Signature(basis, sigWriter, blockSize, sumSize, sigType)
	sigWriter.Flush()
	return signature, sigBuf, err
}

func TestSignature(t *testing.T) {
	tests := []struct {
		hashType string
		sigType  librsync.MagicNumber
		inFile   string
	}{
		{hashType: "md4", sigType: librsync.MD4_SIG_MAGIC,
			inFile: "01"},
		{hashType: "blake2", sigType: librsync.BLAKE2_SIG_MAGIC,
			inFile: "01"},
	}
	for _, tt := range tests {
		name := tt.hashType + ":" + tt.inFile
		inFile := filepath.Join("testdir", "signature.input", tt.inFile+".in")
		expectFile := filepath.Join("testdir", "signature.input", tt.hashType, tt.inFile+".sig")
		var blockSize = uint32(2048)
		var sumSize uint32
		// strong sum sizes used in generating the comparison signatures
		switch tt.hashType {
		case "md4":
			sumSize = uint32(8)
		case "blake2":
			sumSize = uint32(32)
		}
		t.Run(name, func(t *testing.T) {
			t.Logf("Input: %q\n", inFile)
			t.Logf("Comparison: %q\n", expectFile)

			sig, sigBuf, err := sigFromFile(inFile, tt.sigType, sumSize, blockSize)
			if err != nil {
			}
			compareBuf2File(t, sigBuf, expectFile)
			assert.Emptyf(t, sig.Weak2block, "Should not generate hash table on Signature generation")
		})
	}
}

func TestCalcStrongSum(t *testing.T) {
	var (
		data8    = make([]byte, 8)
		data256  = make([]byte, 256)
		data1024 = make([]byte, 1024)
		data2048 = make([]byte, 2048)
	)
	for i := 0; i < 8; i++ {
		data8[i] = byte(i)
	}
	for i := 0; i < 256; i++ {
		data256[i] = byte(i)
	}
	for i := 0; i < 1024; i++ {
		data1024[i] = byte(i)
	}
	for i := 0; i < 2048; i++ {
		data2048[i] = byte(i)
	}
	type args struct {
		data      []byte
		sigType   librsync.MagicNumber
		strongLen uint32
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{name: "md4 8 8", args: args{data: data8, sigType: librsync.MD4_SIG_MAGIC,
			strongLen: 8}, want: []byte{102, 174, 30, 48, 91, 237, 24, 103}},
		{name: "md4 256 8", args: args{data: data256,
			sigType: librsync.MD4_SIG_MAGIC, strongLen: 8}, want: []byte{41,
			138, 5, 188, 80, 110, 30, 205}},
		{name: "md4 1024 8", args: args{data: data1024,
			sigType: librsync.MD4_SIG_MAGIC, strongLen: 8}, want: []byte{90,
			226, 87, 196, 126, 155, 225, 36}},
		{name: "md4 2048 8", args: args{data: data2048,
			sigType: librsync.MD4_SIG_MAGIC, strongLen: 8}, want: []byte{142,
			42, 88, 191, 242, 245, 38, 35}},
		{name: "md4 2048 16", args: args{data: data2048,
			sigType: librsync.MD4_SIG_MAGIC, strongLen: 16}, want: []byte{142,
			42, 88, 191, 242, 245, 38, 35, 121, 210, 146, 54, 27, 35, 227, 129}},
		{name: "blake2 8 8", args: args{data: data8,
			sigType: librsync.BLAKE2_SIG_MAGIC, strongLen: 8},
			want: []byte{119, 6, 93, 37, 182, 34, 168, 37}},
		{name: "blake2 256 8", args: args{data: data256,
			sigType: librsync.BLAKE2_SIG_MAGIC, strongLen: 8},
			want: []byte{57, 167, 235, 159, 237, 193, 154, 171}},
		{name: "blake2 256 16", args: args{data: data256,
			sigType: librsync.BLAKE2_SIG_MAGIC, strongLen: 16},
			want: []byte{57, 167, 235, 159, 237, 193, 154, 171,
				200, 52, 37, 198, 117, 93, 217, 14}},
		{name: "blake2 1024 16", args: args{data: data1024,
			sigType: librsync.BLAKE2_SIG_MAGIC, strongLen: 16},
			want: []byte{241, 85, 31, 238, 178, 82, 199, 230,
				11, 179, 98, 32, 91, 209, 172, 47}},
		{name: "blake2 1024 32", args: args{data: data1024,
			sigType: librsync.BLAKE2_SIG_MAGIC, strongLen: 32},
			want: []byte{241, 85, 31, 238, 178, 82, 199, 230,
				11, 179, 98, 32, 91, 209, 172, 47,
				112, 177, 69, 38, 10, 145, 212, 30,
				140, 93, 10, 24, 117, 73, 165, 242}},
		{name: "blake2 2048 32", args: args{data: data2048,
			sigType: librsync.BLAKE2_SIG_MAGIC, strongLen: 32},
			want: []byte{110, 217, 191, 84, 87, 5, 219, 165,
				151, 30, 131, 161, 242, 164, 106, 157,
				213, 172, 47, 232, 169, 52, 241, 60,
				238, 141, 53, 48, 3, 234, 249, 8}},
		{name: "invalid sigType", args: args{data: data8,
			sigType: librsync.MagicNumber(0), strongLen: 8}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := librsync.CalcStrongSum(tt.args.data, tt.args.sigType, tt.args.strongLen)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalcStrongSum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalcStrongSum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadSignature(t *testing.T) {
	tests := []struct {
		hashName string
		hashType librsync.MagicNumber
		sumSize  uint32
		inFile   string
		wantErr  bool
	}{
		{
			hashName: "md4", hashType: librsync.MD4_SIG_MAGIC,
			sumSize: 8, inFile: "01",
		},
		{
			hashName: "blake2", hashType: librsync.BLAKE2_SIG_MAGIC,
			sumSize: 32, inFile: "01",
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%s %d", tt.hashName, tt.sumSize)
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			inFile := filepath.Join("testdir", "signature.input", tt.hashName, tt.inFile+".sig")
			expectFile := filepath.Join("testdir", "signature.input", tt.inFile+".in")
			var blockSize = uint32(2048)

			expected, _, err := sigFromFile(expectFile, tt.hashType, tt.sumSize, blockSize)
			if err != nil {
				t.Errorf("%v", err)
			}

			basis, err := os.Open(inFile)
			defer basis.Close()
			if err != nil {
				t.Errorf("%v", err)
			}
			got, err := librsync.LoadSignature(basis)
			got.Weak2block = librsync.SignatureHashMap{}

			if err != nil && !tt.wantErr {
				t.Errorf("Unwanted error: %s", err)
			} else if tt.wantErr {
				t.Errorf("Expected error, got none")
			} else {
				assert.Equal(expected, got, "signature")
			}
		})
	}
}

func TestSignatureHashTable(t *testing.T) {
	tests := []struct {
		hashName string
		hashType librsync.MagicNumber
		sumSize  uint32
		inFile   string
		wantErr  bool
	}{
		{
			hashName: "md4", hashType: librsync.MD4_SIG_MAGIC,
			sumSize: 8, inFile: "01",
		},
		{
			hashName: "blake2", hashType: librsync.BLAKE2_SIG_MAGIC,
			sumSize: 32, inFile: "01",
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%s %d", tt.hashName, tt.sumSize)
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			inFile := filepath.Join("testdir", "signature.input", tt.hashName, tt.inFile+".sig")

			basis, err := os.Open(inFile)
			defer basis.Close()
			if err != nil {
				t.Errorf("%v", err)
			}
			loaded, err := librsync.LoadSignature(basis)
			got := loaded.Weak2block

			hm := librsync.NewSignatureHashMap()
			for i, w := range loaded.WeakSigs {
				s := loaded.StrongSigs[i]
				hm.UpdateBlock(w, s, i)
			}

			assert.Equalf(got, hm, "hash table")

			for p, w := range loaded.WeakSigs {
				if h, ok := got.Get(w); !ok {
					assert.Fail("Weak sig not found in hash table, %v", w)
				} else {
					s := loaded.StrongSigs[p]
					if actual, ok := h.Get(s); !ok {
						assert.Fail("Strong sig not found in hash table, %v", s)
					} else {
						assert.Equalf(p, actual, "strong sig file offset in hash-table")
					}
				}
			}
		})
	}
}
