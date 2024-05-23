package asset

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func canonicalName(name string) string {
	return strings.Replace(name, "\\", "/", -1)
}

func filePath(dir, name string) string {
	cName := canonicalName(name)
	return filepath.Join(append([]string{dir}, strings.Split(cName, "/")...)...)
}

func removeSuffixSeparator(name string) string {
	if len(name) > 0 {
		if strings.HasSuffix(name, string(os.PathSeparator)) {
			return name[:len(name)-1]
		}
	}
	return name
}

func addSuffixSeparator(name string) string {
	if !strings.HasSuffix(name, string(os.PathSeparator)) {
		name += string(os.PathSeparator)
	}
	return name
}

func getContent(completePath string, isDir bool) ([]byte, error) {
	if isDir {
		return nil, nil
	}
	return os.ReadFile(completePath)
}

func decodeB64(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, nil
	}
	dst := make([]byte, b64.StdEncoding.DecodedLen(len(src)))
	n, err := b64.StdEncoding.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}

func encodeB64(dec []byte) string {
	return b64.StdEncoding.EncodeToString(dec)
}

func decompressGZip(sDec []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(sDec))
	if err != nil {
		return nil, fmt.Errorf("read %v", err)
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()
	if err != nil {
		return nil, fmt.Errorf("read %v", err)
	}
	if clErr != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func compressGZip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if gz == nil {
		return nil, fmt.Errorf("nil writer")
	}
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func decompressB64GZip(buffer []byte) ([]byte, error) {
	if len(buffer) == 0 {
		return nil, nil
	}
	dec, err := decodeB64(buffer)
	if err != nil {
		return nil, fmt.Errorf("read %v", err)
	}
	return decompressGZip(dec)
}

func compressB64GZip(data []byte) (string, error) {
	if len(data) == 0 {
		return "", nil
	}
	buffer, err := compressGZip(data)
	if err != nil {
		return "", err
	}
	return encodeB64(buffer), nil
}

func decompressGzipWithHash(buffer []byte) ([]byte, string, error) {
	out, err := decompressGZip(buffer)
	if err != nil {
		return nil, "", err
	}
	hash := md5Hash(out)
	return out, hash, nil
}

func md5HashFile(path string) (string, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return md5Hash(buffer), nil
}

func md5Hash(data []byte) string {
	m := md5.Sum(data)
	m = md5.Sum(m[:])
	return hex.EncodeToString(m[:])
}
