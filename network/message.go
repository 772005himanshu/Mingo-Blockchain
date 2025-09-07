package network

import "github.com/772005himanshu/Mingo-Blockchain/core"

type GetBlocksMessage struct {
	From uint32
	// if To is 0 , then we have to return the maximum blocks
	To uint32
}


type BlocksMessage struct {
	Blocks []*core.Block
}

type GetStatusMessage struct {
}

type StatusMessage struct {
	// the id of the server
	ID            string
	CurrentHeight uint32
	Version       uint32
}
