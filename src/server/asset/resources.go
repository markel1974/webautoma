package asset

import (
	"fmt"
)

type Resources struct {
	piContainer map[string]string
	clContainer map[string]IReaderLock
}

func New() *Resources {
	return &Resources{
		piContainer: make(map[string]string),
		clContainer: make(map[string]IReaderLock),
	}
}

func (res *Resources) SetupFiles(pathId map[string]string) error {
	if len(pathId) == 0 {
		return nil
	}
	for path, id := range pathId {
		res.piContainer[path] = id
		res.clContainer[id] = NewReaderLockFile(path, true)
	}
	for _, cl := range res.clContainer {
		if err := cl.Setup(); err != nil {
			return err
		}
	}
	return nil
}

func (res *Resources) SetupBuffer(id string, data []byte) error {
	if len(id) == 0 {
		return nil
	}
	res.clContainer[id] = NewReaderLockBuffer(data, true)

	for _, cl := range res.clContainer {
		if err := cl.Setup(); err != nil {
			return err
		}
	}
	return nil
}

func (res *Resources) Names(id string) []string {
	cl, ok := res.clContainer[id]
	if !ok {
		return nil
	}
	return cl.Names()
}

func (res *Resources) Get(id string, name string) ([]byte, error) {
	cl, ok := res.clContainer[id]
	if !ok {
		return nil, fmt.Errorf("invalid id: %s", id)
	}
	return cl.Get(name)
}

func (res *Resources) GetHash(id string, name string) ([]byte, string, error) {
	cl, ok := res.clContainer[id]
	if !ok {
		return nil, "", fmt.Errorf("invalid id: %s", id)
	}
	return cl.GetHash(name)
}

func (res *Resources) update(path string) error {
	id, ok := res.piContainer[path]
	if !ok {
		return fmt.Errorf("unknown path %s", path)
	}

	cl, ok := res.clContainer[id]
	if !ok {
		return fmt.Errorf("unknown path:id => %s:%s", path, id)
	}

	return cl.Reload()
}
