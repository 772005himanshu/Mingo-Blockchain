package core 

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type Blockchain struct {
	store Storage  // this storage would contains complete blocks of the transactions
	headers []*Header // list of the slice if points to headers , we make the list in the memeory cheap and easy to retrive through it -> Ram is cheap
	validator Validator
}

func NewBlockchain(genesis *Block) (*Blockchain, error) {
	bc := &Blockchain {
		headers: []*Header{},
		store : NewMemoryStore(),
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
	return bc.headers[height] , nil
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

// [0,1,2,3] -> 4 len
// [0,1,2,3] -> 3 Height
func (bc *Blockchain) Height() uint32 {
	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {

	bc.headers = append(bc.headers, b.Header)

	logrus.WithFields(logrus.Fields{
		"height": b.Height,
		"hash": b.Hash(BlockHasher{}),
	}) // whats is this used for ?
	return bc.store.Put(b)
}


/// for the out dated package  - keep updated
/// -> use command : go get -u golang.org/x/sys
