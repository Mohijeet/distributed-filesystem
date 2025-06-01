package main

import (
	"fmt"
	"log"

	"github.com/mohijeet/distributed-filesystem/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
}

type FileServer struct {
	FileServerOpts
	store  *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {

	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
	}
}
func (s *FileServer) Stop() {
	close(s.quitch)
}
func (s *FileServer) loop() {
	defer func() {
		log.Print("stopping server")
		s.Transport.Close()
	}()
	for {
		select {
		case <-s.quitch:
			return
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)

		}
	}
}
func (s *FileServer) Start() error {

	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	s.loop()
	return nil

}
