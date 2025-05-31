package main

import (
	"bytes"
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
		Filename: hashStr,
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
	return PathKey{
		PathName: key,
		Filename: key,
	}
}

func NewStore(opts StoreOpts) *Store {
	if opts.pathTransformFunc == nil {
		opts.pathTransformFunc = DefaultPathTransformFunc
	}
	return &Store{
		StoreOpts: opts,
	}
}

type PathKey struct {
	PathName string
	Filename string
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.Filename)
}

func (p PathKey) ReturnRootPath() string {
	f := strings.Split(p.FullPath(), "/")

	return f[0]
}
func (s *Store) Delete(key string) error {

	pathKey := s.pathTransformFunc(key)

	fmt.Printf("file path be deleted %v\n", pathKey.FullPath())
	defer func() {
		fmt.Printf("Deleted file %s\n", pathKey.Filename)
	}()
	//return os.RemoveAll(pathKey.FullPath())
	return os.RemoveAll(pathKey.ReturnRootPath())

}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)

	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)
	defer f.Close()

	_, err = io.Copy(buff, f)

	return buff, err

}
func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.pathTransformFunc(key)

	return os.Open(pathKey.FullPath())
	// if err != nil {
	// 	return nil, err
	// }
	// // r := new(bytes.Buffer)
	// // io.Copy(r, f)

}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.pathTransformFunc(key)

	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	pathAndFileName := pathKey.FullPath()

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
