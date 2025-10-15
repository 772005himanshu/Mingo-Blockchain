package core

import (
	"fmt"
	"math/rand"

	"github.com/772005himanshu/Mingo-Blockchain/crypto"
	"github.com/772005himanshu/Mingo-Blockchain/types"
)

// We create the collection owner of the collection is the public key of main transaction(parent key)
// The sender of the transaction

type TxType byte

const (
    TxTypeCollection TxType = iota // 0x0 Collection of the NFT we have minted
    TxTypeMint                    // 0x1
)


type CollectionTx struct {
    Fee int64  // if you want to place the Collection we have to pay the fee
    Metadata []byte
}

type MintTx struct {
    Fee int64
    NFT types.Hash // It is the simple byte data holding ut represent the Value
    Collection type.Hash // To separate the different NFT to different types Simply saying
    Metadata      []byte  // Show some data like name , symbol , value it holds
    CollectionOwner crypto.PublicKey
    Signature       crypto.Signature
}

// Only the Public Value to be encoded in the Transactions
type Transaction struct {
    Type TxType
    TxInner any // interface{}
	Data []byte // What is this data mean this is the ByteCode for the VM here 
	From      crypto.PublicKey  // the sender of the transaction
	Signature *crypto.Signature
    Nonce     uint64

	// cached version of the tx data hash
	hash types.Hash
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
		Nonce: rand.Int63n(1000000000), // updating the state every time in the Contract
	}
}

func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() {
		tx.hash = hasher.Hash(tx)
	}
	return tx.hash
}

func (tx *Transaction) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(tx.Data)
	if err != nil {
		return err
	}

	tx.From = privKey.PublicKey()
	tx.Signature = sig

	return nil
}

func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("tx has no Signature")
	}

	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}

func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx)
}

func (tx *Transaction) Encode(enc Encoder[*Transaction]) error {
	return enc.Encode(tx)
}

func init() {
    gob.Register(CollectionTx{})
    gob.Register(MintTx{})
}
