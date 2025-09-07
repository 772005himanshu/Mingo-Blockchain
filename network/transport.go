package network

import "net"

type NetAddr string

type Transport interface {
	Consume() <-chan RPC // taking message that sent to transport layers
	Connect(Transport) error
	SendMessage(net.Addr, []byte) error // Message Shoould be in Addr , bbytes of message
	Broadcast([]byte) error
	Addr() net.Addr
}
