package network

import (
	"fmt"
	"time"
	"github.com/772005himanshu/Mingo-Blockchain/crypto"
	"github.com/772005himanshu/Mingo-Blockchain/core"
	"github.com/sirupsen/logrus"
)

type ServerOpts struct {
	Transports []Transport
	BlockTime time.Duration
	PrivateKey *crypto.PrivateKey
} // this can be used as the Block explorer and wallet and passage to the tx go through it 

type Server struct { // This Behaves as the validator and they also participate in the consensus
	ServerOpts
	blockTime  time.Duration
	memPool *TxPool
	isValidator bool
	rpcCh  chan RPC
	quitCh chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		blockTime : opts.BlockTime,
		memPool: NewTxPool(),
		isValidator: opts.PrivateKey != nil, // this simply means that if you have the private key then you are validtor , if not you are not the validator 
		rpcCh:      make(chan RPC),
		quitCh:     make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.blockTime)

free:
	for {
		select {
		case rpc := <-s.rpcCh: // self rpc
			fmt.Printf("%+v\n", rpc)
		case <-s.quitCh: // do i need to quit the rpc channel
			break free
		// default: // if we could not quit , we got stucked Here
		case <-ticker.C:  // its time to create a new Block
			if s.isValidator{
				s.createNewBlock()
			}
			
		}
	}

	fmt.Println("Server shutdown")
}

func (s *Server) handleTransaction(tx *core.Transaction) error {
	if err := tx.Verify(); err != nil {
		return err
	}

	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"hash" : hash,
		}).Info("transaction already in the mempool")

		return nil
	}

	logrus.WithFields(logrus.Fields{
		"hash" : hash,
	}).Info("adding new tx to the mempool")
	
	return s.memPool.Add(tx)
}

func (s *Server) createNewBlock() error {
	fmt.Println("creating a new Block")
	return nil
}
func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) { // go routine
			for rpc := range tr.Consume() {
				// handle
				s.rpcCh <- rpc
			}
		}(tr)
	}
}
