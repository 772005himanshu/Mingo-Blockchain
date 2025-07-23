package core

import (
	"encoding/binary"
	"io"
	"github.com/772005himanshu/Mingo-Blockchain/types"
)

type Header struct {
	Version uint32 
	PrevBlock types.Hash
	Timestamp uint64
	Height uint32
	
	Nonce uint64 
}

func (h *Header) EncodeBinary(w io.Writer) error{
	if err := binary.Write(w, binary.LittleEndian, &h.Version); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h.PrevBlock); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h.Timestamp); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h.Height); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, &h.Nonce)
} 

// first we encode each type and if we get the block in the byte slice then decode to the Header

func (h *Header) DecodeBinary(r io.Reader) error {

}

type Block struct {
	Header
	Transaction []Transaction
}