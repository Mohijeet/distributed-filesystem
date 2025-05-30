package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func CASPathTransformFunc(key string) PathKey {

	hash := sha1.Sum([]byte(key))

	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr) / blocksize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize

		paths[i] = hashStr[from:to]

	}
	return PathKey{
		PathName: strings.Join(paths, "/"),
		Original: hashStr,
	}

}

type pathTransformFunc func(string) PathKey

type StoreOpts struct {
	pathTransformFunc pathTransformFunc
}

type Store struct {
	StoreOpts
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{}
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

type PathKey struct {
	PathName string
	Original string
}

func (p PathKey) Filename() string {
	return fmt.Sprintf("%s%s", p.PathName, p.Original)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathName := s.pathTransformFunc(key)

	if err := os.MkdirAll(pathName.PathName, os.ModePerm); err != nil {
		return err
	}

	pathAndFileName := pathName.Filename()

	f, err := os.Create(pathAndFileName)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)

	if err != nil {
		return err
	}
	fmt.Printf("copied bytes %v, %s", n, pathAndFileName)
	return nil
}
