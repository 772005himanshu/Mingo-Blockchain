package network

import (
	"bytes"
	"fmt"
	"net"
	"io"
)

type TCPPeer struct {
	conn     net.Conn
	Outgoing bool // if its outgoing , it be true , if its incmoing it to be false
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

func (p *TCPPeer) readLoop(rpcCh chan RPC) {
	buf := make([]byte, 2048) 
	for {
		n, err := p.conn.Read(buf)
		if err == io.EOF {
			continue
		}
		if err != nil {
			fmt.Printf("read error: %s", err)
			continue
		}

		msg := buf[:n]
		rpcCh <- RPC{
			From:    p.conn.RemoteAddr(),
			Payload: bytes.NewReader(msg),
		}

		fmt.Println(string(msg)) // no need
	}
}

type TCPTransport struct {
	peerCh     chan *TCPPeer
	listenAddr string
	listner    net.Listener
}

func NewTCPTransport(addr string, peerCh chan *TCPPeer) *TCPTransport {
	return &TCPTransport{
		peerCh:     peerCh,
		listenAddr: addr,
	}
}

func (t *TCPTransport) Start() error {
	ln, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	t.listner = ln

	go t.acceptLoop()

	return nil

}



func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listner.Accept() // accept from the net library
		// Accept waits for and returns the next connection to the listener.
		if err != nil {
			// we cannot break the loop , if we break it we are break for everyone after that , that is not good for blockchain
			fmt.Printf("accept error from %+v\n", conn)
			continue
		}

		peer := &TCPPeer{
			conn: conn,
		}

		t.peerCh <- peer // here we are passing the peer to channel

	}
}
