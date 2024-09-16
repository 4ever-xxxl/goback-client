package functions

import (
	"bytes"
	"encoding/binary"
)

type LZ77Block struct {
	Offset int
	Length int
	Next   byte
}

func LZ77Compress(data []byte) ([]LZ77Block, error) {
	var blocks []LZ77Block
	windowSize := 4096
	lookaheadSize := 18

	for i := 0; i < len(data); {
		var matchOffset, matchLength int
		for j := 1; j <= windowSize && j <= i; j++ {
			for k := 0; k < lookaheadSize && i+k < len(data); k++ {
				if data[i-j+k] != data[i+k] {
					break
				}
				if k+1 > matchLength {
					matchOffset = j
					matchLength = k + 1
				}
			}
		}

		var next byte
		if i+matchLength < len(data) {
			next = data[i+matchLength]
		}

		blocks = append(blocks, LZ77Block{
			Offset: matchOffset,
			Length: matchLength,
			Next:   next,
		})

		i += matchLength + 1
	}

	return blocks, nil
}

func Compress(inBuf *bytes.Buffer, outBuf *bytes.Buffer) error {
	data := inBuf.Bytes()
	blocks, err := LZ77Compress(data)
	if err != nil {
		return err
	}

	for _, block := range blocks {
		if err := binary.Write(outBuf, binary.LittleEndian, int16(block.Offset)); err != nil {
			return err
		}
		if err := binary.Write(outBuf, binary.LittleEndian, int16(block.Length)); err != nil {
			return err
		}
		if err := outBuf.WriteByte(block.Next); err != nil {
			return err
		}
	}

	return nil
}

func Decompress(inBuf *bytes.Buffer, outBuf *bytes.Buffer) error {
	for {
		var offset int16
		var length int16
		if err := binary.Read(inBuf, binary.LittleEndian, &offset); err != nil {
			break
		}
		if err := binary.Read(inBuf, binary.LittleEndian, &length); err != nil {
			return err
		}
		next, err := inBuf.ReadByte()
		if err != nil {
			return err
		}

		start := outBuf.Len() - int(offset)
		for i := 0; i < int(length); i++ {
			b := outBuf.Bytes()[start+i]
			if err := outBuf.WriteByte(b); err != nil {
				return err
			}
		}

		if err := outBuf.WriteByte(next); err != nil {
			return err
		}
	}

	return nil
}
