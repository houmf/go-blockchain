package main

import (
	"context"
	"fmt"
)

type API struct {
	Blocks map[string]*Block // map[parent_hash]*Block
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
	phs := fmt.Sprintf("%x", block.ParentHash)
	a.Blocks[phs] = block
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
