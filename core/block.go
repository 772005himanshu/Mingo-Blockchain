package core

import (
	"bytes"
	"encoding/gob"
	"io"
	"time"
	"fmt"
	"crypto/sha256"
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

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)

	return buf.Bytes()
}

type Block struct {
	*Header  // it is copied version of the Header -> the * reason behind this we donot maintain the copy of the Header , we want to maintain  a list of the pointers
	Transactions []*Transaction
	Validator crypto.PublicKey
	Signature *crypto.Signature

	// Cached Version of the header Hash 
	hash types.Hash
}

func NewBlock(h *Header,txx []*Transaction) (*Block, error ){
	return &Block{
		Header: h,
		Transactions: txx,
	}, nil
}

func NewBlockFormPrevHeader(prevHeader *Header, txx []*Transaction) (*Block, error) {
	dataHash, err := CalculateDataHash(txx)
	if err != nil {
		return nil, err
	}
	header := &Header{
		Version : 1,
		Height : prevHeader.Height + 1,
		DataHash : dataHash,
		PrevBlockHash: BlockHasher{}.Hash(prevHeader),
		Timestamp: uint64(time.Now().UnixNano()),
	}

	return NewBlock(header,txx)
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)
}


func (b *Block) Sign(privKey crypto.PrivateKey) error  {
	sig , err := privKey.Sign(b.Header.Bytes())
	if err != nil {
		return err // The signature is embedded in the Block then return the error , not the panic
	}

	b.Validator = privKey.PublicKey()
	b.Signature = sig

	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("Block has no signture")
	}

	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("Block has invalid signature")
	}

	for _, tx := range b.Transactions{
		if err := tx.Verify(); err != nil {
			return err
		}
	}

	dataHash, err := CalculateDataHash(b.Transactions)
	if err != nil {
		return err
	}
	if dataHash != b.DataHash {
		return fmt.Errorf("block (%s) has an invalid data hash", b.Hash(BlockHasher{}))
	}

	return nil
}

func (b *Block) Decode(r io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(b)
}

func (b *Block) Encode(w io.Writer, enc Encoder[*Block]) error {
	return enc.Encode(b)
}

func (b *Block) Hash(hasher Hasher[*Header] ) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}

	return b.hash
}

func CalculateDataHash(txx []*Transaction) (hash types.Hash, err error) {
	var (
		buf = &bytes.Buffer{}
	)

	for _, tx := range txx {
		if err = tx.Encode(NewGobTxEncoder(buf)); err != nil {
			return
		}
	}

	hash = sha256.Sum256(buf.Bytes())

	return hash ,err
} 

