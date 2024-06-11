package main

import (
	"encoding/hex"
	"fmt"
	tx "transaction"
)

/*
proof of work

hash256(block header) -> hash value:
0000000000000000007e9e4c586439b0cdbe13b1370bdd9435d76a644d047523

how to check the computed hash is meeting the requirement? => bits e93c0118

from bits field we can compute a valute called target, as long as your hash result
is smaller than the target, then the result can be aceepted

how we can compute target, the target value has characters:
1, it has no more than three none zero value,
2, it has lots of zeros in it
target example:
0000000000000000013ce9000000000000000000000000000000000000000000

let's think how to use 4 bytes to encode the above value?
1, use three bytes to record the none zero bytes in the target
0x18 => 16+8 = 24 bytes
18 01 3c e9 -> reverse byte order -> e8 3c 01 18
e9 3c 01 18 -> little endian formt of above
nonce 4 bytes, miner will guess the value of nonce to make the hash
meet the requirement

e9 3c 01 =reverse byte order => coeefficient: 0x013ce9
last byte in bits : 0x18 , total length of the target in bytes
suffix zeros : 0x18 - 0x03 = exponent: 0x15 , how many zeros at the suffix

left shift 1 place => add 1 bit at the right end with value 0
1 byte => 8 bits => add 1 byte with 0 at the right end, we need to shift 8 places
shift the coefficient to left with (0x15*0x8) places

target value is target = coeffcient << (exponent*8) => coefficient * 2^(exponent * 8)
=> coefficient * (2^8)^exponet => coefficeint * 256^exponent

shift left with 1 place => multiple 2,

2^32 posible value for the nonce, nonce->0, 1, 2....

bitcoin adjustment period, how the difficulty is increase?
when every 2016 blocks are generated, then blockchain will adjust the diffculty
every 2016 blocks are in one group
(1, 2016) (2017, 2017+2016)
1, time_differential = (timestamp of the last block in the last group) - (timestamp of the
first block in current the group) / (2 weeks in seconds)

if time_differential > (2 weeks in seconds) * 4 => time_differential = (2 weeks in seconds) * 4
if time_differential < (2 weeks in seconds) // 4 => time_differential = (2 weeks in seconds)

new_target = previous_target * time_differential // (2 weeks in seconds)

new_target => new_bits
*/

func main() {
	blockRawData, err := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	if err != nil {
		panic(err)
	}

	block := tx.ParseBlock(blockRawData)
	fmt.Printf("block info: %s\n", block)
	blockSerialized := block.Serialize()
	fmt.Printf("serialized block data:%x\n", blockSerialized)

	fmt.Printf("is support BIP0009:%v\n", block.Bip9())
	fmt.Printf("is support BIP0091:%v\n", block.Bip91())
	fmt.Printf("is support BIP0141:%v\n", block.Bip141())

	//256 bits => 32 bytes, => 64 characters in string
	fmt.Printf("target value is %064x\n", block.Target())
	fmt.Printf("difficulty is :%d\n", block.Defficulty().Int64())

	lastBlockRawData, err := hex.DecodeString("00000020fdf740b0e49cf75bb3d5168fb3586f7613dcc5cd89675b0100000000000000002e37b144c0baced07eb7e7b64da916cd3121f2427005551aeb0ec6a6402ac7d7f0e4235954d801187f5da9f5")
	if err != nil {
		panic(err)
	}
	firstBlockRawData, err := hex.DecodeString("000000201ecd89664fd205a37566e694269ed76e425803003628ab010000000000000000bfcade29d080d9aae8fd461254b041805ae442749f2a40100440fc0e3d5868e55019345954d80118a1721b2e")
	if err != nil {
		panic(err)
	}
	newTargt := tx.ComputeNewTarget(firstBlockRawData, lastBlockRawData)
	fmt.Printf("new target is :%064x\n", newTargt.Bytes())

	newBits := tx.TargetToBits(newTargt)
	fmt.Printf("new bits:%x\n", newBits)
}

/*
888 billion
difficulty of first block of bitcoin blockchain to 1

network packet => 1, packet header, 2, payload
block <=> packet => 1, block header, 2, block payload

020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d

1, version, first 4 bytes: 02000020 in little endian

2, the following 32 bytes is the previous block hash:
8ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd000000000000000000
in little endian

3, the following 32 bytes is the merkle root:
5b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be
in little endian

4, following 4 bytes is timestamp :
1e77a759
in little endian format

5, the following 4 bytes call bits, it used in proof-of-work:
e93c0118

6, final 4 bytes is nonce, "number used only once":
a4ffd71d
*/
