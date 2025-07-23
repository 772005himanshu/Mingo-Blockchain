package network

import (
	"fmt"
	"time"
)

type ServerOpts struct {
	Transports []Transport
}

type Server struct {
	ServerOpts
	rpcCh  chan RPC
	quitCh chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		rpcCh:      make(chan RPC),
		quitCh:     make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(5 * time.Second)

free:
	for {
		select {
		case rpc := <-s.rpcCh: // self rpc
			fmt.Printf("%+v\n", rpc)
		case <-s.quitCh: // do i need to quit the rpc channel
			break free
		// default: // if we could not quit , we got stucked Here
		case <-ticker.C:
			fmt.Println("do stuff every x seconds")
		}
	}

	fmt.Println("Server shitdown")
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
