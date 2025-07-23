package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoonnect(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	assert.Equal(t, tra.peers[trb.Addr()], trb)
	assert.Equal(t, trb.peers[tra.Addr()], tra)
}

func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	msg := []byte("Hello world")
	assert.Nil(t, tra.SendMessage(trb.addr, msg))

	rpc := <-trb.Consume() // Comsuming the Channel from the trb to Rpc
	assert.Equal(t, rpc.Payload, msg)
	assert.Equal(t, rpc.From, tra.addr)
}
