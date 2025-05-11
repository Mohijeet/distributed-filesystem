package p2p

// representation of remote node
type Peer interface {
}

// it can handle communication between nodes (tcp, udp, websocket, ...)
type Transport interface {
	ListenAndAccept()
}
