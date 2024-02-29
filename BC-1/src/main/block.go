package main

import (
	"time"
)

type Block struct {
	Time     int64
	Data     []byte
	PrevHash []byte
	Hash     []byte
	Nonce    int
}

func NewBlock(data string, prevHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevHash, []byte{}, 0}
	block.SetHash()
	return block
}

func (b *Block) SetHash() {

	pow := NewProofOfWork(b)
	b.Nonce, b.Hash = pow.Run()
}
