package main

import (
	"time"
)

type Block struct {
	//RlpSerializer
	// The number of this block
	number        uint32
	// Hash to the previous block
	prevHash      string
	// Uncles of this block
	uncles        []*Block
	coinbase      string
	// state xxx
	difficulty    uint32
	// Creation time
	time          int64
	nonce         uint32
	// List of transactions and/or contracts
	transactions  []*Transaction

	extra         string
}

// New block takes a raw encoded string
func NewBlock(raw []byte) *Block {
	block := &Block{}
	block.UnmarshalRlp(raw)

	return block
}

// Creates a new block. This is currently for testing
func CreateBlock(/* TODO use raw data */transactions []*Transaction) *Block {
	block := &Block{
		// Slice of transactions to include in this block
		transactions: transactions,
		number: 1,
		prevHash: "1234",
		coinbase: "me",
		difficulty: 10,
		nonce: 0,

		time: time.Now().Unix(),
	}

	return block
}

func (block *Block) Update() {
}

// Returns a hash of the block
func (block *Block) Hash() string {
	return Sha256Hex(block.MarshalRlp())
}

func (block *Block) MarshalRlp() []byte {
	// Marshal the transactions of this block
	encTx := make([]string, len(block.transactions))
	for i, tx := range block.transactions {
		// Cast it to a string (safe)
		encTx[i] = string(tx.MarshalRlp())
	}

	/* I made up the block. It should probably contain different data or types. It sole purpose now is testing */
	// 測試模擬用
	header := []interface{}{
		block.number,
		block.prevHash,
		//Sha of uncles
		"",
		block.coinbase,
		//root state
		"",
		// Sha of tx
		string(Sha256Bin([]byte(Encode(encTx)))),
		block.difficulty,
		uint64(block.time),
		block.nonce,
		block.extra,
	}

	// TODO
	uncles := []interface{}{}

	// Encode a slice interface which contains the header and the list of transactions.
	return Encode([]interface{}{header, encTx, uncles})
}

func (block *Block) UnmarshalRlp(data []byte) {
	t, _ := Decode(data,0)
	// interface slice assertion
	if slice, ok := t.([]interface{}); ok {
		// interface slice assertion
		if header, ok := slice[0].([]interface{}); ok {
			if number, ok := header[0].(uint8); ok {
				block.number = uint32(number)
			}

			if prevHash, ok := header[1].([]byte); ok {
				block.prevHash = string(prevHash)
			}

			// sha of uncles is header[2]

			if coinbase, ok := header[3].([]byte); ok {
				block.coinbase = string(coinbase)
			}

			// state is header[header[4]

			// sha is header[5]

			// It's either 8bit or 64
			if difficulty, ok := header[6].(uint8); ok {
				block.difficulty = uint32(difficulty)
			}
			if difficulty, ok := header[6].(uint64); ok {
				block.difficulty = uint32(difficulty)
			}

			// It's either 8bit or 64
			if time, ok := header[7].(uint8); ok {
				block.time = int64(time)
			}
			if time, ok := header[7].(uint64); ok {
				block.time = int64(time)
			}

			if nonce, ok := header[8].(uint8); ok {
				block.nonce = uint32(nonce)
			}

			if extra, ok := header[9].([]byte); ok {
				block.extra = string(extra)
			}
		}

		if txSlice, ok := slice[1].([]interface{}); ok {
			// Create transaction slice equal to decoded tx interface slice
			block.transactions = make([]*Transaction, len(txSlice))

			// Unmarshal transactions
			for i, tx := range txSlice {
				if t, ok := tx.([]byte); ok {
					tx := &Transaction{}
					// Use the unmarshaled data to unmarshal the transaction
					// t is still decoded.
					tx.UnmarshalRlp(t)

					block.transactions[i] = tx
				}
			}
		}
	}
}