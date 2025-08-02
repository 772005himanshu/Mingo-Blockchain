package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/772005himanshu/Mingo-Blockchain/crypto"
	"github.com/772005himanshu/Mingo-Blockchain/types"
	"time"
	"fmt"
)


func randomBlock(height uint32) *Block {
	header := &Header {
		Version: 1,
		PrevBlockHash: types.RandomHash(),
		Height: height,
		Timestamp: uint64(time.Now().UnixNano()),
	}

	tx := Transaction{
		Data: []byte("foo"),
	}

	return  NewBlock(header, []Transaction{tx})
}

func TestHashBlock(t *testing.T) {
	b := randomBlock(0)

	fmt.Println(b.Hash(BlockHasher{}))
}

func TestSignBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0)
	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0)
	assert.Nil(t, b.Sign(privKey))
	// assert.NotNil(t, b.Signature)
	assert.Nil(t, b.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()
	assert.NotNil(t, b.Verify())

	b.Height = 100
	assert.NotNil(t, b.Verify())
}