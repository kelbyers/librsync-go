package librsync

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"

	"github.com/resin-os/circbuf"
)

// Delta generates the plan for patching the remote file represented by
// signature `sig`, based on the contents of the local source-of-truth `i`. The
// delta is written to `output`.
func Delta(sig *SignatureType, i io.Reader, output io.Writer) error {
	input := bufio.NewReader(i)

	err := binary.Write(output, binary.BigEndian, DELTA_MAGIC)
	if err != nil {
		return err
	}

	prevByte := byte(0)
	m := match{output: output}

	weakSum := NewRollsum()
	block, _ := circbuf.NewBuffer(int64(sig.BlockLen))
	block.WriteByte(0)
	pos := 0

	for {
		pos++
		in, err := input.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if block.TotalWritten() > 0 {
			prevByte, err = block.Get(0)
			if err != nil {
				return err
			}
		}
		block.WriteByte(in)
		weakSum.Rollin(in)

		if weakSum.count < uint64(sig.BlockLen) {
			continue
		}

		if weakSum.count > uint64(sig.BlockLen) {
			err := m.add(MATCH_KIND_LITERAL, uint64(prevByte), 1)
			if err != nil {
				return err
			}
			weakSum.Rollout(prevByte)
		}

		weakIdx := weakSum.Digest()
		if strongHashes, ok := sig.Weak2block.Get(weakIdx); ok {
			strong2, _ := CalcStrongSum(block.Bytes(), sig.SigType, sig.StrongLen)
			if blockIdx, ok := strongHashes.Get(strong2); ok {
				if bytes.Equal(sig.StrongSigs[blockIdx], strong2) {
					weakSum.Reset()
					block.Reset()
					err := m.add(MATCH_KIND_COPY, uint64(blockIdx)*uint64(sig.BlockLen), uint64(sig.BlockLen))
					if err != nil {
						return err
					}
				}
			}
		}
	}

	for _, b := range block.Bytes() {
		err := m.add(MATCH_KIND_LITERAL, uint64(b), 1)
		if err != nil {
			return err
		}
	}

	if err := m.flush(); err != nil {
		return err
	}

	return binary.Write(output, binary.BigEndian, OP_END)
}
