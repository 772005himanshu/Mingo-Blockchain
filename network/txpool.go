package network

import (
	"github.com/772005himanshu/Mingo-Blockchain/core"
	"github.com/772005himanshu/Mingo-Blockchain/types"
)

type TxPool struct {
	transactions map[types.Hash]*core.Transaction // map with the hash of the transaction and its corresponding Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

// Add an transaction to the pool , the caller is responsible checking if the tx already exist
func (p *TxPool) Add(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	p.transactions[hash] = tx
	return nil
}

func (p *TxPool) Has(hash types.Hash) bool {
	_ , ok := p.transactions[hash]
	return ok
}

func (p *TxPool) Len() int {
	return len(p.transactions)
}


func (p *TxPool) Flush() {
	p.transactions = make(map[types.Hash]*core.Transaction)
} // make a New tx Pool by defining the new map 