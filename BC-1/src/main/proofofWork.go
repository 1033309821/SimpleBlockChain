package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
)

const targetBits = 20

func IntToHex(n int64) []byte {
	// 将 int64 转换为无符号的 byte 数组
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(n))

	// 转换为十六进制字符串
	hexString := make([]byte, 0, 16)
	for _, b := range buf {
		hexString = append(hexString, hexDigits[b>>4])
		hexString = append(hexString, hexDigits[b&0xF])
	}

	// 返回十六进制字节数组
	return hexString
}

// 定义十六进制字符表
var hexDigits = []byte("0123456789abcdef")

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevHash,
			pow.block.Data,
			IntToHex(pow.block.Time),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)

	for {
		here := pow.prepareData(nonce)
		hash = sha256.Sum256(here)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) < 0 {
			break
		}
		nonce++
	}
	fmt.Printf("\r%x", hash)
	fmt.Print("\n\n")
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	hashInt.SetBytes(pow.block.Hash)
	if hashInt.Cmp(pow.target) < 0 {
		return true
	}
	return false
}
