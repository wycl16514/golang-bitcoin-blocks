In previous section, we successfully broadcast our transaction on to the bitcoin blockchain. Since at any given time, there are thousand of transaction waiting for broadcast to chain, every bitcoin node will batch all coming transactions at 
broadcast them at every ten minutes. The collection of  transactions that are broadcast at one time are called block.

We first pay attention to a special transaction that is the first transaction in the block, and it has the name of coinbase transaction. There is a star bitcoin company called coinbase and is already list on Nasdaq, but the coinbase transaction 
we are going to look has nothing to do with it. The coinbase transaction is every important for bitcoin nodes, because nodes can get substantial reward.Let's see an example of coinbase transaction, following is the raw binary data of one coinbase
transaction:
```g
01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff5e03d71b07254d696e656420627920416e74506f6f6c20626a31312f4542312f4144362f43205914293101fabe6d6d678e2c8c34afc36896e7d9402824ed38e856676ee94bfdb0c6c4bcd8b2e5666a0400000000000000c7270000a5e00e00ffffffff01faf20b58000000001976a914338c84849423992471bffb1a54a8d9b1d69dc28a88ac00000000
```
Let's dissect the data block aboved piece by piece:
1, the first four bytes: 01000000 it is the version of the transaction in little endian format
2, the following one byte: 01 is the input count
3, following chunk of zeros: 0000000000000000000000000000000000000000000000000000000000000000, is the previous transaction hash, sinces its the first transaction of the block, therefore it has not previous transaction and this value is all 0,
4, ffffffff previous transaction index
5, the following data chunk: 5e03d71b07254d696e656420627920416e74506f6f6c20626a31312f4542312f4144362f43205914293101fabe6d6d678e2c8c34afc36896e7d9402824ed38e856676ee94bfdb0c6c4bcd8b2e5666a0400000000000000c7270000a5e00e00 is input script
6, ffffffff sequence number
7, 01 output count
8, faf20b5800000000 output amount
9, 1976a914338c84849423992471bffb1a54a8d9b1d69dc28a88ac p2pkh scriptpubkey
10, 00000000 lock time

The structure of coinbase transaction is the same as we have seen before, but has some specials:
1, coinbase transaction must have exactly one input
2, the one input can only have previous transaction id set to 32 bytes of data chunk and filled with all 0
3, the one input can only have previous transaction output index of four bytes with each byte set to value 0xff

Let's add code to check whether a transaction is coinbase or not, in transaction.go add the following code:
```go
func (t *Transaction) IsCoinBase() bool {
	/*
			1, coinbase transaction must have exactly one input
		    2, the one input can only have previous transaction id set to 32 bytes of data chunk
		     and filled with all 0
		    3, the one input can only have previous transaction output index of four bytes
			with each byte set to value 0xff
	*/
	if len(t.txInputs) != 1 {
		return false
	}
	for i := 0; i < len(t.txInputs[0].previousTransaction); i++ {
		if t.txInputs[0].previousTransaction[i] != 0x00 {
			return false
		}
	}

	coinBaseIdx := big.NewInt(int64(0xffffffff))
	if t.txInputs[0].previousTransactionIdex.Cmp(coinBaseIdx) != 0 {
		return false
	}

	return true
}

```
And in main.go we add the following test code:
```go
func main() {
	coinBaseTransactionRawData, err := hex.DecodeString("01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff5e03d71b07254d696e656420627920416e74506f6f6c20626a31312f4542312f4144362f43205914293101fabe6d6d678e2c8c34afc36896e7d9402824ed38e856676ee94bfdb0c6c4bcd8b2e5666a0400000000000000c7270000a5e00e00ffffffff01faf20b58000000001976a914338c84849423992471bffb1a54a8d9b1d69dc28a88ac00000000")
	if err != nil {
		panic(err)
	}
	coinBaseTx := tx.ParseTransaction(coinBaseTransactionRawData)

	fmt.Printf("is coin base transaction: %v\n", coinBaseTx.IsCoinBase())
}
```
Then run the aboved code and we get the following result:
```go
is coin base transaction: true
```
