package main

import (
	"log"

	"github.com/mohijeet/distributed-filesystem/p2p"
)

func main() {
	tr := p2p.NewTCPTransport(":3000")
	err := tr.ListenAndAccept()
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
