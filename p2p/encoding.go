package p2p

import (
	"encoding/gob"
	"fmt"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBdecoder struct{}

func (dec GOBdecoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(r io.Reader, msg *RPC) error {

	buff := make([]byte, 1028)

	n, err := r.Read(buff)
	if err != nil {
		fmt.Printf("response error %v\n", err)
		return err
	}

	msg.PlayLoad = buff[:n]
	return nil

}
