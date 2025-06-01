package p2p

// representation of remote node
type Peer interface {
	Close() error
}

// it can handle communication between nodes (tcp, udp, websocket, ...)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
