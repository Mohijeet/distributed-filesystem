package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddress string
	listener      net.Listener

	shakeHands HandshakeFunc
	mu         sync.RWMutex
	peers      map[net.Addr]Peer
	decoder    Decoder
}

type TCPPeer struct {
	conn net.Conn

	//if we dial connetion and retrive a conn => outbound == true
	//if we accept connection and retrive a conn => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		shakeHands:    NOPHandshakeFunc,
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("tcp accept error %s\n", err)
		}
		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.shakeHands(conn); err != nil {

	}
	msg := &Temp{}
	//read-loop
	for {

		if err := t.decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error %s\n", err)
			continue
		}
	}
	fmt.Printf("new incoming connection %+v\n", peer)
}
