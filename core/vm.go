package core

import (
	"encoding/binary"
)

type Instruction byte

const (
	InstrPushInt Instruction = 0x0a // 10 why not we use 0 , 1, 2 .. 9 they are reserved by the variable 
	InstrAdd Instruction = 0x0b  // 11
	InstrPushByte Instruction = 0x0c // 12
	InstrPack Instruction = 0x0d // 13
	InstrSub Instruction = 0x0e // 14
	InstrStore Instruction = 0x0f  // Store on the Blockchain by serializing
	InstrGet Instruction = 0xae  // getting the data stored on the blockchain
	InstrMul Instruction = 0xea
	InstrDiv Instruction = 0xfd
)

type Stack struct {
	data []any // interface{} type alias 
	sp int  // stack pointer where we are on the stack
}

func NewStack(size int) *Stack {
	return &Stack {
		data: make([]any, size),
		sp: 0,
	}
}

func (s *Stack) Push(v any) {
	s.data = append([]any{v}, s.data...)
	// s.data[s.sp] = v
	s.sp++
}

func (s *Stack) Pop() any {
	value := s.data[0]
	s.data  = append(s.data[:0], s.data[1:]...) // deleting the first in the pre allocated array and shifting the rest of the array by one slot
	s.sp--
	return value
}
type VM struct {
	data []byte
	ip int // instruction pointer
	stack *Stack
	contractState *State
}

func NewVM(data []byte, contractState *State) *VM {
	return &VM{
		contractState: contractState,
		data : data,
		ip : 0,
		stack: NewStack(128),
	}
}

func (vm *VM) Run() error {
	for {
		instr := Instruction(vm.data[vm.ip])

		if err := vm.Exec(instr); err != nil {
			return err
		}

		vm.ip++

		if vm.ip > len(vm.data) - 1 {
			break
		}
	}

	return nil
}


func (vm *VM) Exec(instr Instruction) error {
	switch(instr) {
	case InstrPushInt:
		vm.stack.Push(int(vm.data[vm.ip - 1]))
	case InstrAdd:
		a := vm.stack.Pop().(int) // we have any on the stack so we are giving the hint that data is the int type
		// for the floating point we have to use the Big Int for the bigger space allocation here  
		b := vm.stack.Pop().(int)
		c := a  + b
		vm.stack.Push(c)

	case InstrPushByte:
		vm.stack.Push(byte(vm.data[vm.ip - 1]))
	case InstrPack:
		n := vm.stack.Pop().(int)
		b := make([]byte, n)  // packing the bytes into the string vector 

		for i := 0; i< n ; i++ {
			b[i] = vm.stack.Pop().(byte) // then push them on the string vector 
		}

		vm.stack.Push(b)
	case InstrSub:
		a := vm.stack.Pop().(int) 
		b := vm.stack.Pop().(int)
		c := a  - b
		vm.stack.Push(c)

	case InstrStore:

		var (
			key = vm.stack.Pop().([]byte) // always store the key in the string format
			value = vm.stack.Pop()

			serializeValue []byte
		)

		switch v := value.(type) {
		case int:
			serializeValue = serializeInt64(int64(v))


		default:
			panic("Unknown type")
		}
		vm.contractState.Put(key , serializeValue)

	case InstrGet:
		var (
			key = vm.stack.Pop().([]byte)
		)

		value , err := vm.contractState.Get(key)
		if err != nil {
			return err
		}

		vm.stack.Push(value)


	case InstrMul:
		a := vm.stack.Pop().(int) 
		b := vm.stack.Pop().(int)
		c := a * b
		vm.stack.Push(c)
	

	case InstrDiv:
		b := vm.stack.Pop().(int) 
		a := vm.stack.Pop().(int)
		if (b == 0) {
			break
		}
		c := a / b
		vm.stack.Push(c)
	}

	return nil
}


func serializeInt64(value int64) []byte {
	buf := make([]byte, 8)

	binary.LittleEndian.PutUint64(buf, uint64(value))

	return buf
}

func deserializeInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}