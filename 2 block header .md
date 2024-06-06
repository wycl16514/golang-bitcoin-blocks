Block in blockchain is like packet for the intertnet, the raw data of transaction is like payload of packet for internet. Since internet packet has header, and block also has its own 
header. Let's dissect a block binary data into fields, following is an example of raw data for a block:

```go
020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d
```

1, first 4 bytes is version in little endian, 02000020

2, following 32 bytes of data chunk in little endian indicate the previous block: 8ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd000000000000000000, this is the hash254 result of the
previous block

3, the following 32 bytes in little endian is the mekle root, 5b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be, we will goto this concept later on.

4, following 4 bytes is timestamp in little endian format: 1e77a75

5, following 4 bytes are called bits: e93c0118

5, the final 4 bytes is called nonce: a4ffd71d

Compare with transaction, fileds in block header are in fix length, and the total length of block header is 80 bytes.

Let's create a block object as the transaction object we have done before, create a new file named block.go in transaction and add the following code:

```go
package transaction

import (
	"bufio"
	"bytes"
	ecc "elliptic_curve"
	"fmt"
	"io"
	"math/big"
)

type Block struct {
	version         []byte
	previousBlockID []byte
	merkleRoot      []byte
	timeStamp       []byte
	bits            []byte
	nonce           []byte
}

func ParseBlock(rawBlock []byte) *Block {
	block := &Block{}

	reader := bytes.NewReader(rawBlock)
	bufReader := bufio.NewReader(reader)
	buffer := make([]byte, 4)
	io.ReadFull(bufReader, buffer)
	block.version = reverseByteSlice(buffer)

	buffer = make([]byte, 32)
	io.ReadFull(bufReader, buffer)
	block.previousBlockID = reverseByteSlice(buffer)

	buffer = make([]byte, 32)
	io.ReadFull(bufReader, buffer)
	block.merkleRoot = reverseByteSlice(buffer)

	buffer = make([]byte, 4)
	io.ReadFull(bufReader, buffer)
	block.timeStamp = reverseByteSlice(buffer)

	buffer = make([]byte, 4)
	io.ReadFull(bufReader, buffer)
	block.bits = buffer

	buffer = make([]byte, 4)
	io.ReadFull(bufReader, buffer)
	block.nonce = buffer

	return block
}

func (b *Block) Serialize() []byte {
	result := make([]byte, 0)
	version := new(big.Int)
	version.SetBytes(b.version)
	result = append(result, BigIntToLittleEndian(version, LITTLE_ENDIAN_4_BYTES)...)
	result = append(result, reverseByteSlice(b.previousBlockID)...)
	result = append(result, reverseByteSlice(b.merkleRoot)...)

	timeStamp := new(big.Int)
	timeStamp.SetBytes(b.timeStamp)
	result = append(result, BigIntToLittleEndian(timeStamp, LITTLE_ENDIAN_4_BYTES)...)

	result = append(result, b.bits...)
	result = append(result, b.nonce...)

	return result
}

func (b *Block) Hash() []byte {
	s := b.Serialize()
	sha := ecc.Hash256(string(s))
	return reverseByteSlice(sha)
}

func (b *Block) String() string {
	s := fmt.Sprintf("version:%x\nprevious block id:%x\nmerkle root:%x\ntime stamp:%x\nbits:%x\nnonce:%x\nhash:%x\n",
		b.version, b.previousBlockID, b.merkleRoot, b.timeStamp, b.bits, b.nonce, b.Hash())
	return s
}
```
Let's test aboved code in main.go as following:
```go
package main

import (
	"encoding/hex"
	"fmt"
	tx "transaction"
)

func main() {
blockRawData, err := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	if err != nil {
		panic(err)
	}
	block := tx.ParseBlock(blockRawData)
	fmt.Printf("block info: %s\n", block)
	blockSerialized := block.Serialize()
	fmt.Printf("serialized block data: %x\n", blockSerialized)
}
```
Run the aboved code and we have the following result:
```go
block info: version:20000002
previous block id:000000000000000000fd0c220a0a8c3bc5a7b487e8c8de0dfa2373b12894c38e
merkle root:be258bfd38db61f957315c3f9e9c5e15216857398d50402d5089a8e0fc50075b
time stamp:59a7771e
bits:e93c0118
nonce:a4ffd71d
hash:0000000000000000007e9e4c586439b0cdbe13b1370bdd9435d76a644d047523

serialized block data: 020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d
```

Let's see what those fields used for. The version field used to indicate what kind of feature the block can support, version 2 means the block can support protocol BIP0034,
this protocol introduce the block height as we have seen aboved, version 3 means the block will support protocol BIP0065 which specified command OP_CHECKLOCKTIMEVERIFY. How the bits
set of the version will indicate what kind of protocol the miner will support, if the left most threee bits are 001, then the block or the miner will support protocol BIP0009.

Let's add a method to check the version field and find out what kind of protocol the given miner for the block is supporting:
```go
func (b *Block) Bip9() bool {
	//is the miner support BIP0009
	version := new(big.Int)
	version.SetBytes(b.version)
	ver := version.Int64()
	return (ver >> 29) == 0b001
}

func (b *Block) Bip91() bool {
	//is support BIP0091
	version := new(big.Int)
	version.SetBytes(b.version)
	ver := version.Int64()
	return (ver >> 4 & 1) == 1
}

func (b *Block) Bip141() bool {
	//is support BIP0141
	version := new(big.Int)
	version.SetBytes(b.version)
	ver := version.Int64()
	return (ver >> 1 & 1) == 1
}
```
Then we can call the aboved code in main.go:
```go
fmt.Printf("is support BIP0009: %v\n", block.Bip9())
fmt.Printf("is support BIP0091: %v\n", block.Bip91())
fmt.Printf("is support BIP0141: %v\n", block.Bip141())
```
The result for the aboved code is :
```g
is support BIP0009: true
is support BIP0091: false
is support BIP0141: true
```

All blocks are link together with one by one, that's why we call block chain, given one block, we can traverse back to the first block by using the previous block field, we will dive
deeper to this in later section. Merkle root is used for verification for some purposes, and we will go to it in detail in later sections.Timestamp used to verify lock time for given 
transaction and compute the difficulty for every 2016 blocks, we will goto it at later sections. Bits used to support the so call proof-of-work, we will discuss it later on, nonce is
"number used only once", it also used in proof-of-work.
