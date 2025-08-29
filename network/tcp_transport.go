package network

import (
	"net"
	"fmt"
)

type TCPTransport struct {
	listenAddr string
	listner    net.Listener
}

func NewTCPTransport(addr string) *TCPTransport {
	return &TCPTransport{
		listenAddr: addr,
	}
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn , err := t.listner.Accept()
		if err != nil {
			// we cannot break the loop , if we break it we are break for everyone after that , that is not good for blockchain
			fmt.Printf("accept error from %+v\n", conn)
			continue
		}

		fmt.Printf("%+v\n", conn)
	}
}

func (t *TCPTransport) Start() error {
	ln , err := net.Listen("tcp",t.listenAddr)
	if err != nil {
		return err
	}

	t.listner = ln

	go t.acceptLoop()

	fmt.Println("TCP listen to port: ", t.listenAddr)

	return nil

}