package asset

import (
	"fmt"
	"sync"
)

type ReaderLockFile struct {
	decompress bool
	lock       sync.RWMutex
	path       string
	hash       string
	reader     *Reader
}

func NewReaderLockFile(path string, decompress bool) *ReaderLockFile {
	return &ReaderLockFile{
		decompress: decompress,
		path:       path,
		hash:       "",
		reader:     nil,
	}
}

func (rl *ReaderLockFile) Setup() error {
	return rl.Reload()
}

func (rl *ReaderLockFile) Reload() error {
	hash, err := md5HashFile(rl.path)
	if err != nil {
		return err
	}
	rl.lock.RLock()
	equal := hash == rl.hash
	rl.lock.RUnlock()
	if equal {
		return nil
	}
	reader := NewReader(rl.decompress)
	if err := reader.ReadFile(rl.path, ""); err != nil {
		return fmt.Errorf("error while opening asset file (%s): %s", rl.path, err.Error())
	}
	rl.lock.Lock()
	rl.reader = reader
	rl.hash = hash
	rl.lock.Unlock()

	return nil
}

func (rl *ReaderLockFile) Get(name string) ([]byte, error) {
	var content []byte
	var err error

	rl.lock.RLock()
	if rl.reader == nil {
		err = fmt.Errorf("nil reader")
	} else {
		content, err = rl.reader.Get(name)
	}
	rl.lock.RUnlock()

	return content, err
}

func (rl *ReaderLockFile) GetHash(name string) ([]byte, string, error) {
	var content []byte
	var hash string
	var err error

	rl.lock.RLock()
	if rl.reader == nil {
		err = fmt.Errorf("nil reader")
	} else {
		content, hash, err = rl.reader.GetHash(name)
	}
	rl.lock.RUnlock()

	return content, hash, err
}

func (rl *ReaderLockFile) Names() []string {
	var names []string

	rl.lock.RLock()
	if rl.reader != nil {
		names = rl.reader.Names()
	}
	rl.lock.RUnlock()

	return names
}
