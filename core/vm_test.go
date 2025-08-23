package core

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)


func TestStack(t *testing.T) {
	s := NewStack(128) // size of the stack 

	s.Push(1)
	s.Push(2)

	value := s.Pop()

	assert.Equal(t, value, 1)


	fmt.Println(s.data)

	value = s.Pop()
	assert.Equal(t, value, 2)
}

func TestVM(t *testing.T) {
	data := []byte{0x02 , 0x0a, 0x03, 0x0a, 0x0b, 0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0x0f}  // just the pushing from the front Reverse order that we start from 
	pushFoo := []byte{0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c ,0x03, 0x0a, 0x0d, 0xae}
	// F O O => Pack[F O O]

	data = append(data, pushFoo...)

	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())


	value := vm.stack.Pop().([]byte)
	valueSerialized := deserializeInt64(value)

	assert.Equal(t, valueSerialized, int64(5))

	// valueBytes, err := contractState.Get([]byte("FOO"))
	// value := deserializeInt64(valueBytes)
	// assert.Nil(t, err)
	// assert.Equal(t,value, int(5))

	
}