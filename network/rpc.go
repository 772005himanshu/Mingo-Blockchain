package network

import (
	"bytes"
	"encoding/gob"
	"io"
	"fmt"
	"github.com/772005himanshu/Mingo-Blockchain/core"
	"github.com/sirupsen/logrus"
)

type MessageType byte

const (
	MessageTypeTx MessageType = 0x1
)

type RPC struct {
	From    NetAddr
	Payload io.Reader
}

type Message struct {
	Header MessageType
	Data   []byte
}

func NewMessage(t MessageType , data []byte) *Message {
	return &Message{
		Header: t,
		Data: data,
	}
}

func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}

type DecodedMessage struct {
	From NetAddr
	Data any
}

type RPCDecode func() (*DecodedMessage, error)

// type RPCHandler interface { // Something like decoder
// 	// convert the plane byte payload  into some logic(message)
// 	HandleRPC(rpc RPC) error
// }

func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error) {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil , fmt.Errorf("failed to decode message from %s: %s", rpc.From, err)
	}

	logrus.WithFields(logrus.Fields{
		"from": rpc.From,
		"type": msg.Header,
	}).Debug("new incoming message")

	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}

		return &DecodedMessage {
			From: rpc.From,
			Data: tx,
		}, nil

	default:
		return nil, fmt.Errorf("invalid message header %x", msg.Header)
	}
} 



type RPCProcessor interface {
	// Take the encoded stuff from the handler and process it , there should have method we call like Rust Solana Native match instruction
	ProcessMessage(*DecodedMessage) error
}
