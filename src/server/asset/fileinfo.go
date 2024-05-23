package asset

import (
	"os"
	"time"
)

type fileInfoData struct {
	FileName    string    `json:"name"`
	FileSize    int64     `json:"size"`
	FileMode    int       `json:"mode"`
	FileModTime time.Time `json:"modTime"`
	FileBuffer  string    `json:"buffer"`
	decompress  bool
	buffer      []byte
	hash        string
	err         error
}

func newFileInfoData(fInfo os.FileInfo, data []byte) (*fileInfoData, error) {
	buffer := ""
	if len(data) > 0 {
		var err error
		if buffer, err = compressB64GZip(data); err != nil {
			return nil, err
		}
	}
	fInfoData := &fileInfoData{
		FileName:    fInfo.Name(),
		FileSize:    fInfo.Size(),
		FileMode:    int(fInfo.Mode()),
		FileModTime: fInfo.ModTime(),
		FileBuffer:  buffer,
		decompress:  false,
		buffer:      nil,
		err:         nil,
	}
	return fInfoData, nil
}

func (fi *fileInfoData) Setup(decompress bool) error {
	fi.decompress = decompress
	if fi.IsDir() {
		return nil
	}
	if fi.buffer, fi.err = decodeB64([]byte(fi.FileBuffer)); fi.err != nil {
		return fi.err
	}
	fi.FileBuffer = ""
	if !fi.decompress {
		return nil
	}
	if fi.buffer, fi.hash, fi.err = decompressGzipWithHash(fi.buffer); fi.err != nil {
		return fi.err
	}
	return nil
}

func (fi *fileInfoData) Name() string {
	return fi.FileName
}

func (fi *fileInfoData) Size() int64 {
	return fi.FileSize
}

func (fi *fileInfoData) Mode() os.FileMode {
	return os.FileMode(fi.FileMode)
}

func (fi *fileInfoData) ModTime() time.Time {
	return fi.FileModTime
}

func (fi *fileInfoData) IsDir() bool {
	return os.FileMode(fi.FileMode)&os.ModeDir != 0
}

func (fi *fileInfoData) Sys() interface{} {
	return nil
}

func (fi *fileInfoData) Get() ([]byte, error) {
	if fi.IsDir() {
		return nil, nil
	}
	if !fi.decompress {
		return decompressGZip(fi.buffer)
	}
	return fi.buffer, fi.err
}

func (fi *fileInfoData) GetHash() ([]byte, string, error) {
	if fi.IsDir() {
		return nil, "", nil
	}
	if !fi.decompress {
		return decompressGzipWithHash(fi.buffer)
	}
	return fi.buffer, fi.hash, fi.err
}
