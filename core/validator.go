package core

import (
	"fmt"
)

// Validator Construct the Block and propose them to the network 

// why the interface we can mock it for the testing 
type Validator interface {
	ValidateBlock(*Block) error // we donot know what we are validating 
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc : bc,
	}
}

func (v *BlockValidator) ValidateBlock(b *Block) error { 
	// The validator Block height should not be less than the curremt Block height that's mature because that will means that we already have that block
	// we also cannot implement the block with the height greater then the current height + 1
	if v.bc.HasBlock(b.Height) {
		return fmt.Errorf("chain already contains block (%d) with hash (%s) ", b.Height, b.Hash(BlockHasher{}))
	}

	if b.Height != v.bc.Height() + 1 {
		return fmt.Errorf("block (%s) to high", b.Hash(BlockHasher{}))
	}

	prevHeader , err := v.bc.GetHeader(b.Height - 1) // to Check the prev Hash os the block really exists or not?
	if err != nil {
		return err
	}

	hash := BlockHasher{}.Hash(prevHeader)

	if hash != b.PrevBlockHash{
		return fmt.Errorf("the hash of the previous block (%s) is invalid", b.PrevBlockHash)
	}

	if err := b.Verify(); err != nil {
		return err
	}
	return nil
}