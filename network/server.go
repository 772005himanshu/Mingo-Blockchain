package network

import (
	"fmt"
	"time"
	"bytes"
	"github.com/772005himanshu/Mingo-Blockchain/crypto"
	"github.com/772005himanshu/Mingo-Blockchain/core"
	"github.com/772005himanshu/Mingo-Blockchain/types"
	"github.com/sirupsen/logrus"
	"github.com/go-kit/log"
	"os"
)


var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	ID string
	Logger log.Logger
	RPCDecodeFunc   RPCDecodeFunc
	RPCProcessor RPCProcessor
 	Transports []Transport
	BlockTime time.Duration
	PrivateKey *crypto.PrivateKey
} // this can be used as the Block explorer and wallet and passage to the tx go through it 

type Server struct { // This Behaves as the validator and they also participate in the consensus
	ServerOpts
	memPool *TxPool
	chain  *core.Blockchain
	isValidator bool
	rpcCh  chan RPC
	quitCh chan struct{}
}

func NewServer(opts ServerOpts) (*Server ,error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime 
	}

	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}

	chain, err := core.NewBlockchain(genesisBlock())
	if err != nil {
		return nil , err
	}

	s := &Server{
		ServerOpts: opts,
		chain: chain,
		memPool: NewTxPool(),
		isValidator: opts.PrivateKey != nil, // this simply means that if you have the private key then you are validtor , if not you are not the validator 
		rpcCh:      make(chan RPC),
		quitCh:     make(chan struct{}, 1),
	}

	// If we donot got any processor from the server options , we going to user
	// the server as default
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	return s, nil

}

func (s *Server) Start() {
	s.initTransports()

free:
	for {
		select {
		case rpc := <-s.rpcCh: // self rpc
		    msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				s.Logger.Log("error", err)
			}

			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				s.Logger.Log("error", err)
			}
		case <-s.quitCh: // do i need to quit the rpc channel
			break free
		// default: // if we could not quit , we got stucked Here
		// ISSUE - If the server donot have the validator , we donot want to create this ticker
		// it will never get hit because of the upper ProcessMessage Take time to solve and verify it is correct or not ?
		// case <-ticker.C:  // its take time to create a new Block
		// 	if s.isValidator{
		// 		s.createNewBlock()
		// 	}
			
		// }
		}

	}

	s.Logger.Log("msg", "Server is Shutting down")
}



func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.BlockTime)

	s.Logger.Log("msg", "Starting validator loop", "blockTime", s.BlockTime)

	for {
		<-ticker.C
		s.createNewBlock()
	}
}

func (s *Server) ProcessMessage(msg *DecodedMessage) error {

	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	}

	return nil
}


func (s *Server) broadcast(msg []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(payload) ;  err != nil {
			return err
		}
	}

	return nil
} 

func (s *Server) processTransaction(tx *core.Transaction) error {

	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	logrus.WithFields(logrus.Fields{
		"hash" : hash,
		"mempool length": s.memPool.Len(),
	}).Info("adding new tx to the mempool")

	s.Logger.Log("msg", "adding new tx to mempool",
		"hash", hash,
		"mempoolLength", s.memPool.Len(),
	)

	// TODO : broadcast this tx to peers
	go s.broadcastTx(tx)
	
	return s.memPool.Add(tx)
}


func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())
}

func (s *Server) createNewBlock() error {
	currentHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}

	block, nil := core.NewBlockFormPrevHeader(currentHeader, nil)
	if err != nil {
		return err
	}

	if err := block.Sign(*s.PrivateKey); err != nil {
		return err
	}

	if err := s.chain.AddBlock(block); err != nil {
		return err
	}

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


func genesisBlock() *core.Block {
	header := &core.Header {
		Version : 1,
		DataHash : types.Hash{},
		Height: 0,
		Timestamp: uint64(time.Now().UnixNano()),
	}

	b, _ := core.NewBlock(header, nil)
	return b
}