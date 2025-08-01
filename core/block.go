package core

import (
	"bytes"
	"encoding/gob"
	"io"
	"fmt"

	"github.com/772005himanshu/Mingo-Blockchain/crypto"

	"github.com/772005himanshu/Mingo-Blockchain/types"
)

type Header struct {
	Version   uint32
	DataHash types.Hash
	PrevBlockHash types.Hash
	Timestamp uint64
	Height    uint32
	Nonce     uint64
}



type Block struct {
	*Header  // it is copied version of the Header -> the * reason behind this we donot maintain the copy of the Header , we want to maintain  a list of the pointers
	Transactions []Transaction
	Validator crypto.PublicKey
	Signature *crypto.Signature

	// Cached Version of the header Hash 
	hash types.Hash
}

func NewBlock(h *Header,tx []Transaction) *Block {
	return &Block{
		Header: h,
		Transactions: tx,
	}
}


func (b *Block) Sign(privKey crypto.PrivateKey) *crypto.Signature{
	sig , err := privKey.Sign(b.HeaderData())
	if err != nil {
		return nil // The signature is embedded in the Block then return the error , not the panic
	}

	b.Validator = privKey.PublicKey()
	b.Signature = sig

	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("Block has no signture")
	}

	if !b.Signature.Verify(b.Validator, b.HeaderData()) {
		return fmt.Errorf("Block has invalid signature")
	}

	return nil
}

func (b *Block) Decode(r io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(r, b)
}

func (b *Block) Encode(w io.Writer, enc Encoder[*Block]) error {
	return enc.Encode(w, b)
}

func (b *Block) Hash(hasher Hasher[*Block] ) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b)
	}

	return b.hash
}

func (b *Block) HeaderData() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(b.Header)

	return buf.Bytes()
}
