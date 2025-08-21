package core

import (
	"io"
	"crypto/elliptic"
	"encoding/gob"
)

// if we place the encoding in the network , then we have the circular dependencies 

type Encoder[T any] interface {
	Encode(T) error
}

type Decoder[T any] interface {
	Decode(T) error
}

type GobTxEncoder struct {
	w io.Writer
}


func NewGobTxEncoder(w io.Writer) *GobTxEncoder {
	return &GobTxEncoder {
		w: w,
	}
}

func (e *GobTxEncoder) Encode(tx *Transaction) error {
	enc := gob.NewEncoder(e.w)  // The Encoder is responsible to include a buffer or connection or whatever si it can stream its encoding
	return enc.Encode(tx)
}

type GobTxDecoder struct {
	r io.Reader
}

func NewGobTxDecoder(r io.Reader) *GobTxDecoder {
	gob.Register(elliptic.P256())
	return &GobTxDecoder {
		r: r,
	}
}

func (e *GobTxDecoder) Decode(tx *Transaction) error {
	enc := gob.NewDecoder(e.r)  
	return enc.Decode(tx)
}

type GobBlockEncoder struct {
	w io.Writer
}

func NewGobBlockEncoder(w io.Writer) *GobBlockEncoder {
	return &GobBlockEncoder {
		w: w,
	}
}

func (enc *GobBlockEncoder) Encode(b *Block) error {
	return gob.NewEncoder(enc.w).Encode(b)
}

type GobBlockDecoder struct {
	r io.Reader
}

func NewGobBlockDecoder(r io.Reader) *GobBlockDecoder {
	return &GobBlockDecoder {
		r: r,
	}
}

func (dec *GobBlockDecoder) Decode(b *Block) error {
	return gob.NewDecoder(dec.r).Decode(b)
}

func init() {
	gob.Register(elliptic.P256())
}