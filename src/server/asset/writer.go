package asset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Writer struct {
	data map[string]*fileInfoData
}

func NewWriter() *Writer {
	return &Writer{
		data: make(map[string]*fileInfoData),
	}
}

func (h *Writer) Setup(path string) error {
	if len(path) == 0 {
		return fmt.Errorf("empty path")
	}
	path = canonicalName(path)
	//path = addSuffixSeparator(path)
	path = removeSuffixSeparator(path)
	return h.directoryWalk(path, path)
}

func (h *Writer) Marshal(packageName string, varName string, fileName string) error {
	if len(packageName) == 0 {
		packageName = "asset"
	}
	if len(varName) == 0 {
		varName = "Content"
	}

	content, err := json.Marshal(h.data)
	if err != nil {
		return err
	}
	compressedContent, err := compressB64GZip(content)
	if err != nil {
		return err
	}
	if err = os.WriteFile(fileName, []byte(compressedContent), 0644); err != nil {
		return err
	}
	f, err := os.OpenFile(fileName+".go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write([]byte("package " + packageName + "\n\n")); err != nil {
		return err
	}
	if _, err := f.Write([]byte("const " + varName + "=`")); err != nil {
		return err
	}
	if _, err := f.Write([]byte(compressedContent)); err != nil {
		return err
	}
	if _, err := f.Write([]byte("`")); err != nil {
		return err
	}
	return nil
}

func (h *Writer) directoryWalk(base string, pathReal string) error {
	files, err := ioutil.ReadDir(pathReal)
	if err != nil {
		return err
	}

	for _, fInfo := range files {
		var data []byte
		pathRealChild := pathReal + string(os.PathSeparator) + fInfo.Name()
		if !strings.HasPrefix(pathRealChild, base) {
			continue
		}
		if fInfo.IsDir() {
			if err := h.directoryWalk(base, pathRealChild); err != nil {
				return err
			}
		}
		data, err := getContent(pathRealChild, fInfo.IsDir())
		if err != nil {
			return err
		}
		fInfoData, err := newFileInfoData(fInfo, data)
		if err != nil {
			return err
		}
		target := pathRealChild[len(base)+1:]
		h.data[target] = fInfoData
	}
	return nil
}
