package core

import (
	"testing"
	"github.com/772005himanshu/Mingo-Blockchain/crypto"
	"github.com/stretchr/testify/assert"
	"bytes"
)

func TestSignTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	data := []byte("foo")
	tx := &Transaction {
		Data: data,
	}

	assert.Nil(t, tx.Sign(privKey)) // Sign the transaction and we verify there is no error here 
	assert.NotNil(t, tx.Signature ) // make sure the signture is not null
}

func TestVerifyTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t,tx.Sign(privKey))
	assert.Nil(t,tx.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	tx.From = otherPrivKey.PublicKey()

	assert.NotNil(t, tx.Verify())
}


func TestTxEncodeDecode(t *testing.T) {
	tx := randomTxWithSignature(t)
	buf := &bytes.Buffer{}
	assert.Nil(t, tx.Encode(NewGobTxEncoder(buf)))

	txDecoded := new(Transaction)
	assert.Nil(t, txDecoded.Decode(NewGobTxDecoder(buf)))
	assert.Equal(t, tx, txDecoded)
} 

func randomTxWithSignature(t *testing.T) *Transaction {
	privKey := crypto.GeneratePrivateKey()
	tx := Transaction{
		Data : []byte("foo"),	
	}
	assert.Nil(t, tx.Sign(privKey))

	return &tx
}