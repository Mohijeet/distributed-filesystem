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

const defaultRootFolderName = "DFS"

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

type PathTransformFunc func(string) PathKey

type StoreOpts struct {
	Root              string
	PathTransformFunc PathTransformFunc
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
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
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
func (s *Store) StateFile(key string) bool {
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, s.PathTransformFunc(key).FullPath())
	fmt.Printf("file existance checking... %s", fullPathWithRoot)
	_, err := os.Stat(fullPathWithRoot)

	// if os.IsNotExist(err) {
	// 	return false
	// }
	return os.IsNotExist(err)
	//return !errors.Is(err, os.ErrNotExist)
}
func (s *Store) Delete(key string) error {

	// err := s.StateFie(key)
	// if err != nil {
	// 	return err
	// }
	pathKey := s.PathTransformFunc(key)

	fmt.Printf("file path be deleted %v\n", pathKey.FullPath())
	defer func() {
		fmt.Printf("Deleted file %s\n", pathKey.Filename)
	}()

	//return os.RemoveAll(pathKey.FullPath())
	fullNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.ReturnRootPath())
	return os.RemoveAll(fullNameWithRoot)

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
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	return os.Open(fullPathWithRoot)
	// if err != nil {
	// 	return nil, err
	// }
	// // r := new(bytes.Buffer)
	// // io.Copy(r, f)

}

func (s *Store) Write(key string, r io.Reader) error {
	return s.writeStream(key, r)

}
func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.PathName)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())

	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)

	if err != nil {
		return err
	}
	fmt.Printf("copied bytes %v, %s", n, fullPathWithRoot)

	return nil
}
