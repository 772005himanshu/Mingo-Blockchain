package core 

import (
	"fmt"
	"sync"
	"github.com/go-kit/log"
	
)

type Blockchain struct {
	logger log.Logger
	// Distributed System this need to thread safe -> we can add the mutex or construct something mechanism (avoiding  the mutex by using the channels)
	store Storage  // this storage would contains complete blocks of the transactions
	lock sync.RWMutex
	headers []*Header // list of the slice if points to headers , we make the list in the memeory cheap and easy to retrive through it -> Ram is cheap
	validator Validator
}

func NewBlockchain(l log.Logger, genesis *Block) (*Blockchain, error) {
	bc := &Blockchain {
		headers: []*Header{},
		store : NewMemoryStore(),
		logger: l,
	}

	bc.validator = NewBlockValidator(bc) // validator should be constructed from the config file
	// so every default interface implementation should be constructed from the config file so people wnat to replace it they just swap it out in the config 

	err := bc.addBlockWithoutValidation(genesis)
	return bc, err
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) AddBlock(b *Block) error {
	// validate 
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	return bc.addBlockWithoutValidation(b)
}


func (bc *Blockchain) GetHeader(height uint32) (*Header, error ) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high" , height)
	}
	bc.lock.Lock() // what is the use of this lock here
	defer bc.lock.Unlock() // then on next there is unlock
	return bc.headers[height] , nil  // We are going to grab the header from the list 
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

// [0,1,2,3] -> 4 len
// [0,1,2,3] -> 3 Height
func (bc *Blockchain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	
	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.lock.Unlock()

	bc.logger.Log(
		"msg", "new block",
		"hash" , b.Hash(BlockHasher{}),
		"height", b.Height,
		"transactions", len(b.Transactions),
	)
	
	return bc.store.Put(b) // put the block in the Storage 
}


/// for the out dated package  - keep updated
/// -> use command : go get -u golang.org/x/sys
