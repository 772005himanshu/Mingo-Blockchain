package network 

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/772005himanshu/Mingo-Blockchain/core"
	"strconv"
	"math/rand"
)


func TestTxPool(t *testing.T) {
	p := NewTxPool()
	assert.Equal(t, p.Len(), 0)
}


func TestTxPoolAddTx(t *testing.T) {
	p := NewTxPool()
	tx := core.NewTransaction([]byte("foo"))
	assert.Nil(t, p.Add(tx))
	assert.Equal(t, p.Len(), 1)

	_ = core.NewTransaction([]byte("foo"))
	assert.Equal(t, p.Len(), 1)

	p.Flush()
	assert.Equal(t, p.Len(), 0)
}


func TestSortTransactions(t *testing.T) {
	p := NewTxPool()
	txLen := 1000

	for i := 0 ; i< txLen; i++ {
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i) , 10)))  // if we pass the same tx again and again this should be count as the one(Becasue the Hash of tx is the same ) in the TX POOL
		tx.SetFirstSeen(int64(i * rand.Intn(1000)))  // Intn is in the Math library
		assert.Nil(t, p.Add(tx))
	} // If we are passing the 1000 tx in the txPool there should be 1000 tx in the TX Pool

	assert.Equal(t, txLen, p.Len())

	txx := p.transactions
	for i := 0; i< len(txx) - 1; i++ {
		assert.True(t, txx[i].FirstSeen() < txx[i + 1].FirstSeen())
	}
	
}