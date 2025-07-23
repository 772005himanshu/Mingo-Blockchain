package main

import (
	"github.com/772005himanshu/Mingo-Blockchain/network"
	"time"
)

// Server  -> container
// Transport  -> tcp, udp
// Block
// Tx

func main() {
	trLocal := network.NewLocalTransport("LOCAL")   // Your node
	trRemote := network.NewLocalTransport("REMOTE") // Some Body Else Node

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for { // Remote Node Sending Message Every second
			trRemote.SendMessage(trLocal.Addr(), []byte("Hello Mingo Blockchain"))
			time.Sleep(1 * time.Second)
		}
	}()

	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal},
	}

	s := network.NewServer(opts)
	s.Start()
}
