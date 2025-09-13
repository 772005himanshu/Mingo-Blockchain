package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/772005himanshu/Mingo-Blockchain/core"
	"github.com/772005himanshu/Mingo-Blockchain/crypto"
	"github.com/772005himanshu/Mingo-Blockchain/types"
	"github.com/go-kit/log"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	SeedNodes     []string
	ListenAddr    string
	TCPTransport  *TCPTransport
	ID            string
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
}

type Server struct {
	TCPTransport *TCPTransport
	peerCh       chan *TCPPeer         // to communicate with the peers
	peerMap      map[net.Addr]*TCPPeer // we have track the peer in the server not in the tcp Layer -> this help in boardcasting it to the each node or validator by looping through all the peer s
	ServerOpts
	mempool     *TxPool
	chain       *core.Blockchain
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
	mu          sync.RWMutex
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "addr", opts.ID)
	}

	chain, err := core.NewBlockchain(opts.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}

	peerCh := make(chan *TCPPeer)
	tr := NewTCPTransport(opts.ListenAddr, peerCh)

	s := &Server{
		TCPTransport: tr,
		// We have two type of channel -> blocking and infinite channel
		peerCh:      peerCh,                      // we need to add the blocking channel why ? we need to implement the deterministic way to handle the blockchain and with out the race condition (reentrancy Peer repaeting like that )
		peerMap:     make(map[net.Addr]*TCPPeer), //
		ServerOpts:  opts,
		chain:       chain,
		mempool:     NewTxPool(1000),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}

	s.TCPTransport.peerCh = peerCh // this is specific path we are using

	// If we dont got any processor from the server options, we going to use
	// the server as default.
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	// for _, tr := range s.Transports {
	// 	if err := s.sendGetStatusMessage(tr); err != nil {
	// 		s.Logger.Log("send get status error", err)
	// 	}
	// }

	// s.bootstrapNodes()

	return s, nil
}

func (s *Server) bootstrapNetwork() {
	// This functionality basic use is to loop through all the Nodes -> Connection
	// pass it the PeerCh
	for _, addr := range s.SeedNodes {
		fmt.Println("trying to connect to", addr)
		go func(addr string) {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Printf("could not connect to %+v\n", conn)
				return
			}
			s.peerCh <- &TCPPeer{
				conn: conn,
			}
		}(addr)
	}
}

func (s *Server) Start() {
	s.TCPTransport.Start()

	time.Sleep(time.Second * 1)

	s.bootstrapNetwork()

	s.Logger.Log("accepting TCP connection on", "addr", s.ListenAddr, "id", s.ID)

free:
	for {
		select {
		case peer := <-s.peerCh:
			// TODO - add the Mutex next time
			s.peerMap[peer.conn.RemoteAddr()] = peer
			go peer.readLoop(s.rpcCh) // fix this

			if err := s.sendGetStatusMessage(peer); err != nil {
				s.Logger.Log("err", err)
				continue
			}

			s.Logger.Log("msg", "peer added to the server", "outgoing", peer.Outgoing , "addr", peer.conn.RemoteAddr())

		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				s.Logger.Log("error", err)

				continue
			}

			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				if err != core.ErrBlockKnown {
					s.Logger.Log("error", err)
				}
			}

		case <-s.quitCh:
			break free
		}
	}

	s.Logger.Log("msg", "Server is shutting down")
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
	case *core.Block:
		return s.processBlock(t)

	case *GetStatusMessage:
		return s.processGetStatusMessage(msg.From, t)

	case *StatusMessage:
		return s.processStatusMessage(msg.From, t)

	case *GetBlocksMessage:
		return s.processGetBlocksMessage(msg.From, t)

	case *BlocksMessage:
		return s.processBlocksMessage(msg.From, t)
	}
	return nil
}

func (s *Server) processGetBlocksMessage(from net.Addr, data *GetBlocksMessage) error {
	s.Logger.Log("msg" ,"received getBlocks message", "from", from)

	var (
		blocks = []*core.Block{}

		height = s.chain.Height()
	)

	if data.To == 0 {
		for i := int(data.From);i< int(height); i++ {
			block, err := s.chain.GetBlock(uint32(i))
			if err != nil {
				return err
			}

			blocks = append(blocks, block)
		}
	}

	blocksMsg := &BlocksMessage{
		Blocks: blocks,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(blocksMsg); err != nil {
		return err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	msg := NewMessage(MessageTypeGetBlocks, buf.Bytes())

	peer, ok := s.peerMap[from]
	if !ok {
		return fmt.Errorf("peer %s ot known", peer.conn.RemoteAddr())
	}

	return peer.Send(msg.Bytes())

}

// TODO: Remove the logic from the main function to here
// Normally Transport which is our transport should do the trick

func (s *Server) sendGetStatusMessage(peer *TCPPeer) error {
	var (
		getStatusMsg = new(GetStatusMessage)
		buf          = new(bytes.Buffer)
	)

	if err := gob.NewEncoder(buf).Encode(getStatusMsg); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeGetStatus, buf.Bytes())

	return peer.Send(msg.Bytes())




}

func (s *Server) broadcast(payload []byte) error {
	// We have to lock this here, accesssing the peer and adding to the peers -> its is a basically a concurrent data race condition problem @note
	// So why we need the sync Mutex -> so all the peers to be sync
	s.mu.RLock()
	defer s.mu.RUnlock()
	for netAddr, peer := range s.peerMap {
		if err := peer.Send(payload); err != nil {
			fmt.Printf("peer send error => addr %s [err : %s]\n", netAddr, err)
		}
	}
	return nil
}

func (s *Server) processBlocksMessage(from net.Addr, data *BlocksMessage) error {
	s.Logger.Log("msg", "received Blocks", "from", from)

	for _, block := range data.Blocks {
		fmt.Printf("BLOCK => %+v\n", block)
		if err := s.chain.AddBlock(block); err != nil {
			fmt.Printf("Adding block Error %s\n", err)
			continue
		}
	}
	return nil
}

func (s *Server) processStatusMessage(from net.Addr,data *StatusMessage) error {
	s.Logger.Log("msg", "received STATUS message", "from", from) // why we are using the Logger , because we have the prefix to it , and we will wait what node is sending what

	if data.CurrentHeight <= s.chain.Height() {
		s.Logger.Log("msg", "cannot sync blockHeight to low","our Height", s.chain.Height(),"their height", data.CurrentHeight, "addr", from)
		return nil
	}

	// In this case we are 100% sure that the node has blocks heigher than us
	getBlockMessage := &GetBlocksMessage {
		From : s.chain.Height(),
		To : 0,
	}
	buf := new(bytes.Buffer)

	if err := gob.NewEncoder(buf).Encode(getBlockMessage); err != nil {
		return err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	peer, ok := s.peerMap[from]
	if !ok {
		return fmt.Errorf("peer %s ot known", peer.conn.RemoteAddr())
	}
	msg := NewMessage(MessageTypeGetBlocks, buf.Bytes())

	return peer.Send(msg.Bytes())
}

func (s *Server) processGetStatusMessage(from net.Addr, data *GetStatusMessage) error {
	s.Logger.Log("msg", "received getStatus message", "from", from) // why we are using the Logger , because we have the prefix to it , and we will wait what node is sending what

	statusMessage := &StatusMessage{
		CurrentHeight: s.chain.Height() + 1,
		ID:            s.ID,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(statusMessage); err != nil {
		return err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	peer, ok := s.peerMap[from]
	if !ok {
		return fmt.Errorf("peer %s ot known", peer.conn.RemoteAddr())
	}
	msg := NewMessage(MessageTypeStatus, buf.Bytes())

	return peer.Send(msg.Bytes())
}

func (s *Server) processBlock(b *core.Block) error {
	if err := s.chain.AddBlock(b); err != nil {
		return err
	}

	go s.broadcastBlock(b)

	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.mempool.Contains(hash) {
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	// s.Logger.Log(
	// 	"msg", "adding new tx to mempool",
	// 	"hash", hash,
	// 	"mempoolPending", s.mempool.PendingCount(),
	// )

	go s.broadcastTx(tx)

	s.mempool.Add(tx)

	return nil
}


// TODO : Find a way to amke sure we donot keep syncing when we are at the highest block height in the network
func (s *Server) requestBlocksLoop(peer net.Addr) error {
	ticker := time.NewTicker(3 * time.Second)

	for {
		// Logic

		height := s.chain.Height()

		s.Logger.Log("msg", "requesting new blocks", "currentHeight", height)
		// In this We are 100 % sure that the node has Blocks Heigher than us
		getBlockMessage := &GetBlocksMessage {
			From : height,
			To : 0,
		}
		buf := new(bytes.Buffer)
	
		if err := gob.NewEncoder(buf).Encode(getBlockMessage); err != nil {
			return err
		}
	
		s.mu.RLock()
		defer s.mu.RUnlock()

	    msg := NewMessage(MessageTypeGetBlocks, buf.Bytes())
		peer, ok := s.peerMap[peer]
		if !ok {
			return fmt.Errorf("peer %s ot known", peer.conn.RemoteAddr())
		}

		if err := peer.Send(msg.Bytes()) ; err != nil {
			s.Logger.Log("error", "failed to send to peer", "err", err, "peer", peer)
		}

		<- ticker.C
	    // But This thinker Sometimes Go out of the Sync
		// After every 3 seconds it go and take message through loop
	}
}

func (s *Server) broadcastBlock(b *core.Block) error {
	buf := &bytes.Buffer{}
	if err := b.Encode(core.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeBlock, buf.Bytes())

	return s.broadcast(msg.Bytes())
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

	// For now we are going to use all transactions that are in the pending pool
	// Later on when we know the internal structure of our transaction
	// we will implement some kind of complexity function to determine how
	// many transactions can be included in a block.
	txx := s.mempool.Pending()

	block, err := core.NewBlockFromPrevHeader(currentHeader, txx)
	if err != nil {
		return err
	}

	if err := block.Sign(*s.PrivateKey); err != nil {
		return err
	}

	if err := s.chain.AddBlock(block); err != nil {
		return err
	}

	// TODO(@anthdm): pending pool of tx should only reflect on validator nodes.
	// Right now "normal nodes" does not have their pending pool cleared.
	s.mempool.ClearPending()

	go s.broadcastBlock(block)

	return nil
}

func genesisBlock() *core.Block {
	header := &core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Height:    0,
		Timestamp: 000000,
	}

	b, _ := core.NewBlock(header, nil)

	privKey := crypto.GeneratePrivateKey()
	if err := b.Sign(privKey); err != nil {
		panic(err)
	}
	return b
}

// How we are going to sync the nodes
// When we boot up the node the seeds node -> that means that the  boot node is going to connect with the Seed nodes like that the idea behind this
// checking the version of the node we are going to connect with only connect with node having the same version
