package network

import (
	"bytes"
	"encoding/gob"
	"io"
	"fmt"
	"github.com/772005himanshu/Mingo-Blockchain/core"
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

type RPCHandler interface { // Something like decoder
	// convert the plane byte payload  into some logic(message)
	HandleRPC(rpc RPC) error
}

type DefaultRPCHandler struct {
	p RPCProcessor
}

func NewDefaultRPCHandler(p RPCProcessor) *DefaultRPCHandler {
	return &DefaultRPCHandler{
		p: p,
	}
}

func (h *DefaultRPCHandler) HandleRPC(rpc RPC) error {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return fmt.Errorf("failed to decode message from %s: %s", rpc.From, err)
	}

	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return err
		}

		return h.p.ProcessTransaction(rpc.From, tx)

	default:
		return fmt.Errorf("invalid message header %x", msg.Header)
	}

}

type RPCProcessor interface {
	// Take the encoded stuff from the handler and process it , there should have method we call like Rust Solana Native match instruction
	ProcessTransaction(NetAddr, *core.Transaction) error
}
