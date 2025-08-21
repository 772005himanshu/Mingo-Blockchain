package core

import (
	"fmt"

	"github.com/772005himanshu/Mingo-Blockchain/crypto"
	"github.com/772005himanshu/Mingo-Blockchain/types"
)

// Only the Public Value to be encoded in the Transactions
type Transaction struct {
	Data []byte

	From crypto.PublicKey
	Signature *crypto.Signature


	// cached version of the tx data hash
	hash types.Hash

	// // first Seen is the Timestamp of when this tx is first seen locally 
	// firstSeen int64
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction {
		Data: data,
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

// func (tx *Transaction) SetFirstSeen(t int64) {
// 	tx.firstSeen = t
// }


// func (tx *Transaction) FirstSeen() int64 {
// 	return tx.firstSeen
// }