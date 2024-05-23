package asset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Reader struct {
	data       map[string]*fileInfoData
	tree       map[string][]*fileInfoData
	decompress bool
}

func NewReader(decompress bool) *Reader {
	return &Reader{
		decompress: decompress,
		data:       make(map[string]*fileInfoData),
		tree:       make(map[string][]*fileInfoData),
	}
}

func (h *Reader) ReadFile(filename string, prefix string) error {
	compressedContent, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.ReadBuffer(compressedContent, prefix)
}

func (h *Reader) ReadBuffer(compressedContent []byte, prefix string) error {
	content, err := decompressB64GZip(compressedContent)
	if err != nil {
		return err
	}
	prefix = removeSuffixSeparator(prefix)
	data := make(map[string]*fileInfoData)
	if err = json.Unmarshal(content, &data); err != nil {
		return err
	}
	h.data = make(map[string]*fileInfoData, len(data))
	for path, child := range data {
		if err := h.add(prefix, path, child); err != nil {
			return err
		}
	}
	return nil
}

func (h *Reader) add(prefix string, path string, child *fileInfoData) error {
	if err := child.Setup(h.decompress); err != nil {
		return err
	}
	canonicalName := canonicalName(prefix + path)
	h.data[canonicalName] = child

	dir := filepath.Dir(canonicalName)
	if children, ok := h.tree[dir]; ok {
		children = append(children, child)
		h.tree[dir] = children
	} else {
		h.tree[dir] = []*fileInfoData{child}
	}
	return nil
}

func (h *Reader) Add(realPath string, assetPath string) error {
	fInfo, err := os.Stat(realPath)
	if err != nil {
		return err
	}
	data, err := getContent(realPath, fInfo.IsDir())
	if err != nil {
		return err
	}
	child, err := newFileInfoData(fInfo, data)
	if err != nil {
		return err
	}
	return h.add("", assetPath, child)
}

func (h *Reader) Get(name string) ([]byte, error) {
	canonicalName := canonicalName(name)
	if f, ok := h.data[canonicalName]; ok {
		a, err := f.Get()
		if err != nil {
			return nil, fmt.Errorf("asset %s read error: %v", name, err)
		}
		return a, nil
	}
	return nil, fmt.Errorf("asset %s not found", name)
}

func (h *Reader) GetHash(name string) ([]byte, string, error) {
	canonicalName := canonicalName(name)
	if f, ok := h.data[canonicalName]; ok {
		a, hash, err := f.GetHash()
		if err != nil {
			return nil, "", fmt.Errorf("asset %s read error: %v", name, err)
		}
		return a, hash, nil
	}
	return nil, "", fmt.Errorf("asset %s not found", name)
}

func (h *Reader) Info(name string) (os.FileInfo, error) {
	canonicalName := canonicalName(name)
	if f, ok := h.data[canonicalName]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

func (h *Reader) Names() []string {
	names := make([]string, 0, len(h.data))
	for name := range h.data {
		names = append(names, name)
	}
	return names
}

func (h *Reader) Dir(dir string) ([]string, error) {
	dir = addSuffixSeparator(dir)
	canonicalName := canonicalName(dir)
	path := filepath.Dir(canonicalName)
	children, ok := h.tree[path]
	if !ok {
		return nil, fmt.Errorf("empty name")
	}
	rv := make([]string, 0, len(children))
	for _, child := range children {
		rv = append(rv, child.Name())
	}
	return rv, nil
}

func (h *Reader) Restore(dir, name string) error {
	data, err := h.Get(name)
	if err != nil {
		return err
	}
	info, err := h.Info(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

func (h *Reader) RestoreDir(dir string, name string) error {
	children, err := h.Dir(name)
	if err != nil {
		return h.Restore(dir, name)
	}
	for _, child := range children {
		err = h.Restore(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}
