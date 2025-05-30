package main

import (
	"fmt"
	"log"

	"github.com/mohijeet/distributed-filesystem/p2p"
)

func OnPeer(peer p2p.Peer) error {
	fmt.Print("ONPEER SUCCESS")
	return nil
}

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)
	err := tr.ListenAndAccept()

	go func() {
		for {
			resp := <-tr.Consume()
			fmt.Printf("printing channel %s\n", resp)
		}

	}()
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
