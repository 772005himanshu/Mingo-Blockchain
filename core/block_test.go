package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/772005himanshu/Mingo-Blockchain/crypto"
	"github.com/772005himanshu/Mingo-Blockchain/types"
	"time"
	"fmt"
)


func randomBlock(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey := crypto.GeneratePrivateKey()
	tx := randomTxWithSignature(t)

	header := &Header {
		Version: 1,
		PrevBlockHash: prevBlockHash,
		Height: height,
		Timestamp: uint64(time.Now().UnixNano()),
	}

	b, err := NewBlock(header, []Transaction{tx})
	assert.Nil(t, err)
	dataHash, err := CalculateDataHash(b.Transactions)
	assert.Nil(t, err)
	b.Header.DataHash = dataHash
	assert.Nil(t, b.Sign(privKey))
	return b
}

func TestHashBlock(t *testing.T, prevBlockHash types.Hash)  {
	b := randomBlock(t,0, prevBlockHash)

	fmt.Println(b.Hash(BlockHasher{})) 
}

func TestSignBlock(t *testing.T, prevBlockHash types.Hash) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(t, 0 , types.Hash{})
	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T, prevBlockHash types.Hash) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(t, 0,types.Hash{})
	assert.Nil(t, b.Sign(privKey))
	// assert.NotNil(t, b.Signature)
	assert.Nil(t, b.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()
	assert.NotNil(t, b.Verify())

	b.Height = 100
	assert.NotNil(t, b.Verify())
}


