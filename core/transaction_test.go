package core

import (
	"testing"
	"github.com/772005himanshu/Mingo-Blockchain/crypto"
	"github.com/stretchr/testify/assert"
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
	tx.PublicKey = otherPrivKey.PublicKey()

	assert.NotNil(t, tx.Verify())
}