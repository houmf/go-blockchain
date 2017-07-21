package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"fmt"

	context "golang.org/x/net/context"
)

type API struct {
	Blocks      map[string]*Block // map[parent_hash]*Block
	LatestBlock *Block
}

func (a *API) Get(ctx context.Context, block *Block) (*Blocks, error) {
	nextBlock := block
	blocks := []*Block{}
	for {
		phs := fmt.Sprintf("%x", nextBlock.ParentHash)
		if cblock, ok := a.Blocks[phs]; ok {
			blocks = append(blocks, cblock)
			nextBlock = cblock
		} else {
			break
		}
	}
	return &Blocks{Blocks: blocks}, nil
}

func (a *API) Announce(ctx context.Context, block *Block) (*Empty, error) {
	// TODO Check if new block is valid
	fmt.Printf("Got new block parent=%x hash=%x nonce=%x\n", block.Hash, block.Nonce, block.ParentHash)
	phs := fmt.Sprintf("%x", block.ParentHash)
	a.Blocks[phs] = block
	a.LatestBlock = block
	return nil, nil
}

func (a *API) blockExists(block *Block) bool {
	for _, b := range a.Blocks {
		if string(b.Hash) == string(block.Hash) {
			return true
		}
	}
	return false
}

func (a *API) Mine() {
	for {
		currentLatestBlock := a.LatestBlock
		fmt.Printf("Trying to mine next block parent=%x\n", currentLatestBlock.Hash)
		pow := GetPow()
		nonce := pow.CalcNonce(currentLatestBlock.Data, 3)
		fmt.Printf("Mined new block nonce=%x\n", nonce)

		data := make([]byte, 64)
		rand.Read(data)

		nb := &Block{
			Data:       data[:],
			Nonce:      nonce,
			ParentHash: currentLatestBlock.Hash,
		}

		var buf bytes.Buffer
		binary.Write(&buf, binary.BigEndian, nb)
		hs := sha1.Sum(buf.Bytes())
		nb.Hash = hs[:]

		// TODO announce block to other peers
		ctx := context.Background()
		a.Announce(ctx, nb)
	}
}
