package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "hello"

	pathname := CASPathTransformFunc(key)
	fmt.Print(pathname)
}
func TestStore(t *testing.T) {
	opts := StoreOpts{
		pathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	data := bytes.NewReader([]byte("some jpg bytes inside files"))
	if err := s.writeStream("somekey", data); err != nil {
		t.Error(err)
	}
}
