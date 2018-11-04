package librsync

import (
	"encoding/binary"
	"fmt"
	"io"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/md4"
)

const (
	// Blake2SumLength is the maximum allowed strongsum size for Blake2
	Blake2SumLength = 32
	// Md4SumLength is the maximum allowed strongsum size for MD4
	Md4SumLength = 16
)

// SignatureType holds the signature for a whole file. This includes the all the
// block sums generated for a file and datastructures for fast matching against
// them.
type SignatureType struct {
	SigType    MagicNumber
	BlockLen   uint32
	StrongLen  uint32
	WeakSigs   []uint32
	StrongSigs [][]byte
	Weak2block SignatureHashMap
}

// CalcStrongSum generates the strong checksum for a block
func CalcStrongSum(data []byte, sigType MagicNumber, strongLen uint32) ([]byte, error) {
	switch sigType {
	case BLAKE2_SIG_MAGIC:
		d := blake2b.Sum256(data)
		return d[:strongLen], nil
	case MD4_SIG_MAGIC:
		d := md4.New()
		d.Write(data)
		return d.Sum(nil)[:strongLen], nil
	}
	return nil, fmt.Errorf("Invalid sigType %#x", sigType)
}

// LoadSignature reads a Signature file and returns a *SignatureType containing
// the loaded Signature.
func LoadSignature(r io.Reader) (*SignatureType, error) {
	signature := &SignatureType{}
	err := binary.Read(r, binary.BigEndian, &signature.SigType)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.BigEndian, &signature.BlockLen)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.BigEndian, &signature.StrongLen)
	if err != nil {
		return nil, err
	}
	signature.Weak2block = NewSignatureHashMap()
	for {
		var (
			weakSig    uint32
			strongSigs []byte
		)
		if err = binary.Read(r, binary.BigEndian, &weakSig); err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		signature.WeakSigs = append(signature.WeakSigs, weakSig)
		strongSigs = make([]byte, signature.StrongLen)
		if err = binary.Read(r, binary.BigEndian, &strongSigs); err != nil {
			return nil, err
		}
		signature.Weak2block.UpdateBlock(weakSig, strongSigs, len(signature.StrongSigs))
		signature.StrongSigs = append(signature.StrongSigs, strongSigs)
	}
	return signature, nil
}

// Signature generates the SignatureType for opened file `infile`, and it to the
// opened signature file `outfile`.
func Signature(input io.Reader, output io.Writer, blockLen, strongLen uint32, sigType MagicNumber) (*SignatureType, error) {
	var maxStrongLen uint32

	switch sigType {
	case BLAKE2_SIG_MAGIC:
		maxStrongLen = Blake2SumLength
	case MD4_SIG_MAGIC:
		maxStrongLen = Md4SumLength
	default:
		return nil, fmt.Errorf("invalid sigType %#x", sigType)
	}

	if strongLen > maxStrongLen {
		return nil, fmt.Errorf("invalid strongLen %d for sigType %#x", strongLen, sigType)
	}

	err := binary.Write(output, binary.BigEndian, sigType)
	if err != nil {
		return nil, err
	}
	err = binary.Write(output, binary.BigEndian, blockLen)
	if err != nil {
		return nil, err
	}
	err = binary.Write(output, binary.BigEndian, strongLen)
	if err != nil {
		return nil, err
	}

	block := make([]byte, blockLen)

	var ret SignatureType
	ret.SigType = sigType
	ret.StrongLen = strongLen
	ret.BlockLen = blockLen

	for {
		n, err := io.ReadFull(input, block)
		if err == io.EOF {
			break
		} else if err == io.ErrUnexpectedEOF {
			err = nil
		} else if err != nil {
			return nil, err
		}
		data := block[:n]

		weak := WeakChecksum(data)
		err = binary.Write(output, binary.BigEndian, weak)
		if err != nil {
			return nil, err
		}

		strong, _ := CalcStrongSum(data, sigType, strongLen)
		output.Write(strong)

		ret.WeakSigs = append(ret.WeakSigs, weak)
		ret.StrongSigs = append(ret.StrongSigs, strong)
	}

	return &ret, nil
}
