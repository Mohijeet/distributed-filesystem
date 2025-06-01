package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

// func TestPathTransformFunc(t *testing.T) {
// 	key := "hello"

//		pathname := CASPathTransformFunc(key)
//		fmt.Print(pathname)
//	}
func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	key := "somekey"
	print("writing data")
	s := NewStore(opts)
	writeString := []byte("REAL DATA THAT NEEDS TO BE WRITTEN")
	data := bytes.NewReader(writeString)
	if err := s.writeStream(key, data); err != nil {
		t.Error(err)
	}

	fmt.Printf("\n reading data\n")
	readData, err := s.readStream(key)
	if err != nil {
		t.Error(err)
	}

	readString, err := io.ReadAll(readData)

	if err != nil {
		t.Error(err)
	}
	if string(readString) != string(writeString) {
		t.Errorf("expected %s, received %s", writeString, readString)

	}
	fmt.Print("Data is as expected : \n", string(readString))

}

func TestDeleteFile(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	key := "somekey"

	s := NewStore(opts)

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}

}

func TestFileExist(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	key := "somekey"

	s := NewStore(opts)

	if ok := s.StateFile(key); !ok {
		t.Error("file does not exists")

	}

}
