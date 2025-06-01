package main

import (
	"fmt"
	"log"
	"time"

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
		//OnPeer:        OnPeer,
	}

	tcpTransport := p2p.NewTCPTransport(tcpOpts)

	// go func() {
	// 	for {
	// 		resp := <-tr.Consume()
	// 		fmt.Printf("printing channel %s\n", resp)
	// 	}

	// }()

	FileServerOpts := FileServerOpts{
		StorageRoot:       "3000_server",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}

	s := NewFileServer(FileServerOpts)
	go func() {
		time.Sleep(time.Second * 2)
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
