package main

import (
	"crypto/rand"
	"flag"
	fmt "fmt"
	"log"
	"net"
	"os"
	"os/signal"

	grpc "google.golang.org/grpc"
)

func main() {
	var nodeHost string
	var nodePort string
	var serverPort string

	flag.StringVar(&serverPort, "serverport", "8886", "Port of the server")
	flag.StringVar(&nodeHost, "nodehost", "localhost", "Hostname of the node to connect")
	flag.StringVar(&nodePort, "nodeport", "", "Port of the node to connect")
	flag.Parse()

	data := make([]byte, 64)
	rand.Read(data)

	originBlock := &Block{
		Hash:       []byte("origin"),
		Nonce:      []byte{},
		ParentHash: []byte("origin"),
		Data:       data[:],
	}

	fmt.Printf("Origin block hash=%x\n", originBlock.Hash)

	api := &API{
		Blocks: map[string]*Block{
			"": originBlock,
		},
		LatestBlock: originBlock,
	}

	go api.Mine()

	go func() {
		lis, err := net.Listen("tcp", "0.0.0.0:"+serverPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer()

		RegisterBlockchainServer(s, api)

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	go func() {
		if nodePort != "" {
			conn, err := grpc.Dial(nodeHost+":"+nodePort, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
		}
	}()

	c := make(chan os.Signal, 1000)
	signal.Notify(c, os.Interrupt)

	func() {
		for _ = range c {
			fmt.Println("Exiting...")
			return
		}
	}()
}
