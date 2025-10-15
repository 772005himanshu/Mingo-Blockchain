package core

import (
	"bytes"
	"github.com/772005himanshu/Mingo-Blockchain/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestNFTTransaction(t *Testing.T) {
    collectionTx := CollectionTx{
        Fee : 200,
        Metadata: []byte("The beginning of the new collection"),
    }
    privKey := crypto.GeneratePrivateKey()
    tx := &Transaction{
        TxType: TxTypeCollection,
        TxInner: collectionTx,
    }

    tx.Sign(privKey)

    buf := new(bytes.Buffer)
    assert.Nil(t, gob.NewEncoder(buf).Encode(tx))  // encoding here

    txDecoded := &Transaction{}
    assert.Nil(t, gob.NewDecoder(buf).Decode(txDecoded)) // decoding here
    assert.Equal(t,tx,txDecoded)
}

func TestSignTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	data := []byte("foo")
	tx := &Transaction{
		Data: data,
	}

	assert.Nil(t, tx.Sign(privKey)) // Sign the transaction and we verify there is no error here
	assert.NotNil(t, tx.Signature)  // make sure the signture is not null
}

func TestVerifyTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.Nil(t, tx.Verify())

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
		Data: []byte("foo"),
	}
	assert.Nil(t, tx.Sign(privKey))

	return &tx
}
