package main

import (
	"context"
	"fmt"
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
	if a.blockExists(block) == true {
		return nil, nil
	}
	// TODO Check if new block is valid
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
		nonce, diff := CalcNonce(currentLatestBlock)
		fmt.Println("Found new block nonce=%x with difficulty=%d", nonce, diff)
		// if the latest block is still the same
		if string(currentLatestBlock.Hash) == string(a.LatestBlock) {
			nb := &Block{
				Data:       []byte{},
				Nonce:      nonce,
				ParentHash: currentLatestBlock.Hash,
			}
			// TODO serialize nb
			// TODO nb.Hash = serialized nb
			// TODO announce block to other peers
		}
	}
}
