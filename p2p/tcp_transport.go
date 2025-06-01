package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}
type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcchn   chan RPC

	// mu    sync.RWMutex
	// peers map[net.Addr]Peer
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

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcchn:           make(chan RPC),
	}
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcchn
}
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()
	log.Printf("\n *** TCP SERVER STARTING ON PORT %s ***\n ", t.ListenAddr)
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("tcp accept error %s\n", err)
		}
		go t.handleConn(conn)
		fmt.Print("starting connection")
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	var err error
	defer func() {
		fmt.Printf("ERROR %s\n", err)
		conn.Close()

	}()

	if err = t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {

		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	fmt.Printf("accepted new connection from %s\n", conn.RemoteAddr())
	rpc := RPC{}

	//read-loop
	for {
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("TCP ERROR %s\n", err)
			continue
		}
		rpc.FROM = conn.RemoteAddr()
		//fmt.Printf("message %s\n ", rpc)

		t.rpcchn <- rpc

	}

}
