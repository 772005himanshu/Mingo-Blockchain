package main

import (
	"bytes"
	"log"
	"fmt"
	"encoding/gob"
	"time"
	"github.com/772005himanshu/Mingo-Blockchain/core"
	"github.com/772005himanshu/Mingo-Blockchain/crypto"
	"github.com/772005himanshu/Mingo-Blockchain/network"
	"github.com/sirupsen/logrus"
)

// Server  -> container
// Transport  -> tcp, udp
// Block
// Tx

// func main() {
// 	trLocal := network.NewLocalTransport("LOCAL")   // Your node
// 	trRemote := network.NewLocalTransport("REMOTE") // Some Body Else Node

// 	trLocal.Connect(trRemote)
// 	trRemote.Connect(trLocal)

// 	go func() {
// 		for { // Remote Node Sending Message Every second
// 			if err := sendTransaction(trRemote, trLocal.Addr()); err != nil {
// 				logrus.Error(err)
// 			}
// 			time.Sleep(1 * time.Second)
// 		}
// 	}()

// 	privKey := crypto.GeneratePrivateKey()
// 	opts := network.ServerOpts{
// 		PrivateKey: &privKey,
// 		ID:         "LOCAL",
// 		Transports: []network.Transport{trLocal},
// 	}

// 	s, err := network.NewServer(opts)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	s.Start()
// }  // OLD Implementation is this

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemoteA := network.NewLocalTransport("REMOTE_A")
	trRemoteB := network.NewLocalTransport("REMOTE_B")
	trRemoteC := network.NewLocalTransport("REMOTE_C")

	

	trLocal.Connect(trRemoteA)
	trRemoteA.Connect(trRemoteB)
	trRemoteB.Connect(trRemoteA)
	trRemoteB.Connect(trRemoteC)
	trRemoteA.Connect(trLocal)

	initRemoteServer([]network.Transport{trRemoteA, trRemoteB, trRemoteC})

	go func() {
		for { // Remote Node Sending Message Every second
			if err := sendTransaction(trRemoteA, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	privKey := crypto.GeneratePrivateKey()
	localServer := makeServer("LOCAL", trLocal, &privKey)
	localServer.Start()

}

func initRemoteServer(trs []network.Transport) {
	for i := 0; i < len(trs) ; i++ {
		id := fmt.Sprintf("REMOTE_%d", i)
		s := makeServer(id, trs[i], nil)
		// Seperate Go routines
		go s.Start()
	}
}

func makeServer(id string, tr network.Transport, privKey *crypto.PrivateKey) *network.Server {
	opts := network.ServerOpts{
		Transport : tr,
		PrivateKey: privKey,
		ID:         id,
		Transports: []network.Transport{tr},
	}

	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	return s
}
func sendGetStatusMessage(tr network.Transport, to network.NetAddr) error {
	var (
		getStatusMsg = new(network.GetStatusMessage)
		buf = new(bytes.Buffer)
	)

	if err := gob.NewEncoder(buf).Encode(getStatusMsg); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeGetStatus, buf.Bytes())

	return tr.SendMessage(to,msg.Bytes())

}


func sendTransaction(tr network.Transport, to network.NetAddr) error {
	privKey := crypto.GeneratePrivateKey()

	tx := core.NewTransaction(contract())
	tx.Sign(privKey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}

func contract() []byte {
	data := []byte{0x02 , 0x0a, 0x03, 0x0a, 0x0b, 0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0x0f}  // just the pushing from the front Reverse order that we start from 
	pushFoo := []byte{0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c ,0x03, 0x0a, 0x0d, 0xae}
	// F O O => Pack[F O O]

	data = append(data, pushFoo...)
	return data
}

