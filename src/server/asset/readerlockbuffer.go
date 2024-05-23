package asset

import (
	"fmt"
	"sync"
)

type ReaderLockBuffer struct {
	decompress bool
	lock       sync.RWMutex
	hash       string
	reader     *Reader
	data       []byte
}

func NewReaderLockBuffer(data []byte, decompress bool) *ReaderLockBuffer {
	return &ReaderLockBuffer{
		decompress: decompress,
		data:       data,
		hash:       "",
		reader:     nil,
	}
}

func (rl *ReaderLockBuffer) Setup() error {
	reader := NewReader(rl.decompress)
	if err := reader.ReadBuffer(rl.data, ""); err != nil {
		return fmt.Errorf("error while opening asset: (%s)", err.Error())
	}

	rl.lock.Lock()
	rl.data = nil
	rl.reader = reader
	rl.lock.Unlock()

	return nil
}

func (rl *ReaderLockBuffer) Reload() error {
	return nil
}

func (rl *ReaderLockBuffer) Get(name string) ([]byte, error) {
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

func (rl *ReaderLockBuffer) GetHash(name string) ([]byte, string, error) {
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

func (rl *ReaderLockBuffer) Names() []string {
	var names []string

	rl.lock.RLock()
	if rl.reader != nil {
		names = rl.reader.Names()
	}
	rl.lock.RUnlock()

	return names
}
