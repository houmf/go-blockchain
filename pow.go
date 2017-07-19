package main

import (
	"crypto/sha256"
	"fmt"
	"hash"
)

type Powork struct {
	Sha     hash.Hash //interface. Will be using sha256
	Current []byte
}

func (p *Powork) Reset() {
	for i := 0; i < 32; i++ {
		p.Current[i] = 0
	}
}

func GetPow() *Powork {
	tmp := []byte{}
	for i := 0; i < 32; i++ {
		tmp = append(tmp, 0)
	}
	return &Powork{Sha: sha256.New(), Current: tmp}
}

func (p *Powork) CalcHash(block []byte, nonce []byte) []byte {
	tmp := append(block, p.Current...)
	p.Sha.Write(tmp)
	return p.Sha.Sum(nil)
}
func IsColission(hash []byte, difficulty int) bool {
	for i := 0; i < difficulty; i++ {
		if hash[i] != 0 {
			return false
		}
	}
	return true
}

//BA SHOULD BE 32 bytes long
func Increment(BA []byte) { //build it inversely
	//l := len(BA)
	if len(BA) == 0 {
		BA = append(BA, 0)
	}
	if (BA[0]) < 255 {
		BA[0]++
		return
	}
	//incd := false
	for i := 1; i < len(BA); i++ {
		if BA[i] == 255 {
			BA[i] = 0
		} else {
			BA[i]++
			return
		}
	}

	BA = append(BA, 1)
}

//find the nonce for a block
func (p *Powork) CalcNonce(block []byte, difficulty int) []byte {
	p.Reset()
	maxIterations := 1000000000000 //yeah a bit rough...
	for i := 0; i < maxIterations; i++ {
		thisHash := p.CalcHash(block, p.Current)
		fmt.Print(i, ":") //fmt.Printf("%x", h.Sum(nil)))
		fmt.Printf("%x\n", thisHash)
		if IsColission(thisHash, difficulty) {
			return thisHash
		}
		Increment(p.Current)
	}
	return p.Current
}

// func (p Pow) CalcHash([]byte) []byte {
// 	return p.Sha.Write([]byte("hello world!\n"))
// }

// func main() {
// 	tmp := GetPow()
// 	//h := sha256.New()
// 	//h.Write([]byte("hello world!\n"))
// 	//fmt.Println("h", h)
// 	//fmt.Printf("%x", h.Sum(nil))
// 	//tmp := []byte{0}
// 	block := []byte("lol")
// 	lala := tmp.CalcNonce(block, 1)
// 	fmt.Println("found nonce", lala)
// }
