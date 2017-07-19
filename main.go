package main

import (
	"flag"
	fmt "fmt"
	"log"
	"net"
	"os"
	"os/signal"

	context "golang.org/x/net/context"

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

	go func() {
		originBlock := &Block{
			Hash:       []byte("origin"),
			Nonce:      0,
			ParentHash: []byte("origin"),
			Data:       []byte{},
		}
		lis, err := net.Listen("tcp", "0.0.0.0:"+serverPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer()

		RegisterBlockchainServer(s, &API{
			Blocks: map[string]*Block{
				"": originBlock,
			},
			LatestBlock: originBlock,
		})

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

			c := NewBlockchainClient(conn)

			c.Announce(context.Background(), &Block{Nonce: 123})
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
