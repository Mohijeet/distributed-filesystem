package p2p

import "net"

type RPC struct {
	FROM     net.Addr
	PlayLoad []byte
}
