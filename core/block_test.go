package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/772005himanshu/Mingo-Blockchain/crypto"
	"github.com/772005himanshu/Mingo-Blockchain/types"
	"time"
	"fmt"
)


func randomBlock(height uint32, prevBlockHash types.Hash) *Block {
	header := &Header {
		Version: 1,
		PrevBlockHash: prevBlockHash,
		Height: height,
		Timestamp: uint64(time.Now().UnixNano()),
	}

	return  NewBlock(header, []Transaction{})
}

func randomBlockWithSignature(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(height, prevBlockHash)
	tx := randomTxWithSignature(t)
	b.AddTransaction(tx)
	assert.Nil(t, b.Sign(privKey))
	return b
}

func TestHashBlock(t *testing.T, prevBlockHash types.Hash)  {
	b := randomBlock(0, prevBlockHash)

	fmt.Println(b.Hash(BlockHasher{})) 
}

func TestSignBlock(t *testing.T, prevBlockHash types.Hash) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0 , types.Hash{})
	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T, prevBlockHash types.Hash) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0,types.Hash{})
	assert.Nil(t, b.Sign(privKey))
	// assert.NotNil(t, b.Signature)
	assert.Nil(t, b.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()
	assert.NotNil(t, b.Verify())

	b.Height = 100
	assert.NotNil(t, b.Verify())
}