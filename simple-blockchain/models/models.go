package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

type Block struct {
	Position  int
	Data      BookCheckout
	TimeStamp string
	Hash      string
	PrevHash  string
}

type BookCheckout struct {
	BookID       string `json:"book_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

type Book struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	PublicationDate string `json:"publication_date"`
	ISBN            string `json:"isbn"`
}

type Blockchain struct {
	Blocks []*Block
}

func GenesisBlock() *Block {
	return CreateBlock(&Block{}, BookCheckout{IsGenesis: true})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{GenesisBlock()}}
}

func validBlock(block, prevBlock *Block) bool {
	if prevBlock.Hash != block.PrevHash {
		return false
	}
	if !block.validateHash(block.Hash) {
		return false
	}
	if prevBlock.Position+1 != block.Position {
		return false
	}
	return true
}

func CreateBlock(prevBlock *Block, checkoutItem BookCheckout) *Block {
	block := &Block{}
	block.Position = prevBlock.Position + 1
	block.PrevHash = prevBlock.Hash
	block.TimeStamp = time.Now().String()
	block.generateHash()

	return block
}

func (b *Block) generateHash() {
	bytes, _ := json.Marshal(b.Data)

	data := string(b.Position) + b.TimeStamp + string(bytes) + b.PrevHash

	hash := sha256.New()
	hash.Write([]byte(data))

	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

func (b *Block) validateHash(hash string) bool {
	b.generateHash()
	return b.Hash == hash
}

func (bc *Blockchain) AddBlock(data BookCheckout) {

	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	block := CreateBlock(prevBlock, data)

	if validBlock(block, prevBlock) {
		bc.Blocks = append(bc.Blocks, block)
	}
}