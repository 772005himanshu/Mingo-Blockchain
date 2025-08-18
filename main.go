package main

import (
	"github.com/772005himanshu/Mingo-Blockchain/network"
	"github.com/772005himanshu/Mingo-Blockchain/core"
	"time"
	"strconv"
	"math/rand"
	"github.com/772005himanshu/Mingo-Blockchain/crypto"
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/go-kit/log"
	
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
			if err := sendTransaction(trRemote, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	privKey := crypto.GeneratePrivateKey()
	opts := network.ServerOpts{
		PrivateKey : &privKey,
		ID : "LOCAL",
		Transports: []network.Transport{trLocal},
	}

	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}
	s.Start()
}


func makeServer(id string, tr network.Transport, privKey *crypto.PrivateKey) *network.Server {
	opts := network.ServerOpts{
		PrivateKey: privKey,
		ID: id,
		Transports: []network.Transport(tr),
	}
}

func sendTransaction(tr network.Transport, to network.NetAddr) error  {
	privKey := crypto.GeneratePrivateKey()
	data := []byte(strconv.FormatInt(int64(rand.Intn(1000)), 10))
	tx := core.NewTransaction(data)
	tx.Sign(privKey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)) ; err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to , msg.Bytes())
}