// Code generated by go-bindata. DO NOT EDIT.
// sources:
// enum.tmpl (7.552kB)

package assets

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _enumTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x59\x4f\x8f\xdb\xb8\x15\x3f\x5b\x9f\xe2\xad\x90\xec\x4a\xae\x23\xa7\x68\xd1\x43\x16\x73\x58\x64\xd3\x60\x17\xcd\x24\xc0\xa4\x7b\x09\x82\x80\x96\x9e\x6c\xee\x48\xa4\x96\xa4\x3c\x72\x55\x7d\xf7\xe2\x91\x94\x2c\x6b\xe4\x99\x69\x76\x06\xb9\x18\xb2\xf8\xf8\xfe\xff\xf9\x91\x6a\xdb\x17\x90\x61\xce\x05\x42\xb8\x43\x96\xa1\x0a\xbb\x2e\x58\xaf\xe1\xb5\xcc\x10\xb6\x28\x50\x31\x83\x19\x6c\x0e\xb0\x95\x2f\x50\xd4\x25\x2d\xfe\xfc\x1e\x2e\xdf\x7f\x84\x37\x3f\xff\xf2\xf1\xbb\x20\xa8\x58\x7a\xcd\xb6\x08\x6d\x9b\xf8\xc7\xae\x0b\x02\x5e\x56\x52\x19\x88\x02\x00\x80\x30\x2f\x4d\x18\xc4\x41\xdb\xa2\xc8\xe0\x05\xad\x8f\x25\x13\x5f\x92\x9b\x4a\xa1\x69\x0b\xad\x3d\xa3\x97\x97\xac\x44\x78\x75\x01\x09\xfd\x49\xec\x3f\xda\x6c\xd7\xf7\x4c\x69\x5a\xcb\x78\x6a\x20\x2c\x98\x36\x32\xcf\x35\x9a\x10\x5e\x7a\x22\x50\x4c\x6c\x11\x9e\xa9\x5f\x44\x86\xcd\x8a\xb6\x14\xf5\x88\xdf\x6f\xf4\x57\x43\xd7\x05\x0b\xcb\x91\x78\xbc\xb7\x3c\x88\xa6\x2a\xea\xf4\xfa\x94\xb1\x93\xf9\x5f\xc8\xb9\xd2\x06\xba\xae\x6d\xe1\x99\x1c\x36\xe8\x7a\xe3\x45\x38\xce\xbd\x60\x2f\x00\x78\x0e\xf8\x47\x4f\x61\x6d\x09\xbf\x84\x5d\xb7\x5e\xc3\xd5\x35\xaf\x2a\xcc\xc0\x2e\xb5\x2d\x16\x1a\xed\xfb\xb6\xf5\xd4\x1f\x14\xe6\xbc\xc1\x8c\x76\x75\x1d\x70\x0d\x8c\x16\x7b\x17\x75\x1d\xc8\x1c\xcc\xa1\xc2\xe3\x16\xf7\xde\x3a\xbc\x37\x90\xe7\xbd\xf4\xd7\xb2\x2c\x51\x18\x5a\x18\x8b\x19\xbd\x26\x7a\xb7\x95\xe2\x77\x4e\x91\xa3\x55\xde\xd4\x97\xd6\x2b\x63\xc5\x2e\x80\x4b\xc3\x1c\xa1\x40\x78\x39\x78\xac\xeb\xe0\x2f\x30\xf2\xe0\xa0\xac\x73\x80\xa7\x1f\x07\x65\x4c\x79\x5b\xc4\x59\x6e\xcf\xbe\xd8\xe8\x10\x03\x1b\xbf\xd3\x90\xba\x07\x9f\x54\xce\xe2\x98\xb2\x13\x0c\x96\x55\xc1\x0c\x42\xa8\x8d\xe2\x62\x8b\x2a\x84\x84\x62\x49\x15\xf0\x81\x29\x8d\x6d\x7b\xcc\xcb\xae\x03\x66\x68\x8b\xd1\x60\x24\xa4\x52\xec\x51\x19\x60\xe0\x36\xd3\x3b\x0a\xd9\x78\x43\x90\xd7\x22\x9d\xe3\x14\x09\x4a\x0e\xb7\x31\x86\xe8\x74\x71\x05\xa8\x94\x54\x31\xb4\xc1\x82\xe7\xd0\xac\x40\x5e\x93\x7d\x5f\x4e\xc9\x6c\x06\x7e\x22\x46\x9f\x7f\x24\x8a\x36\x58\x2c\x14\x9a\x5a\x09\xda\x22\x78\x11\x2c\xba\xb6\xe5\x39\x24\x42\xa6\x4c\x23\xf8\x5c\x78\x4d\xcf\x5c\x68\x14\x9a\x1b\xbe\x47\xa8\x48\xbf\x15\x64\xa4\xbf\xc6\x8a\x51\x3f\x80\x42\xca\xeb\xba\x22\xa3\x2a\x85\x7b\x14\x06\x6a\x21\x30\x45\xad\x99\x3a\x40\x2a\xb5\xa1\x84\x2c\xe4\x0d\xaa\x94\x69\xb2\x7f\x70\x04\xcf\xe1\x06\x21\x93\xe2\x07\x03\x02\x31\x03\x23\x93\x07\x58\xe2\x76\xeb\xe4\xa3\xfc\x17\x71\xb5\x2e\x8a\xef\x32\xad\x0f\xe6\xc2\x5b\xc9\x4a\xd4\xb6\x2f\xf4\xb4\x13\x9f\xbf\x8c\x57\x90\x97\x26\x79\x43\xde\xcd\xa3\xf0\xb9\xa6\x32\x13\x92\x62\xb8\x67\x05\xcf\x60\x1a\x07\xa3\x0e\xf0\xe9\xb9\xfe\x1c\xae\x80\xb8\xaf\xa0\xd7\xf1\x57\xc9\x45\x34\xb1\x82\x7e\xf5\x0a\xc2\x15\x84\x71\xec\xcb\x8b\xb2\xfc\x11\x35\xf2\x7a\xc4\xe3\xe2\xb5\x8d\x96\x7c\x9e\x54\x46\x41\x9f\x73\x51\x33\xd9\x1a\xc3\x07\xa3\xa2\x18\x96\x93\x94\x6e\x07\xdd\xbe\x6f\x82\x2e\xe8\xbb\x49\xcf\xb3\x64\x4a\xef\x58\x01\x6e\x64\xbc\x73\xff\x3e\x62\x63\x80\x97\x55\x81\xd4\x4b\x34\x98\x1d\x82\xa1\x77\x9e\xba\x40\x05\x25\x9a\x9d\xcc\xce\x2a\x33\xe2\x14\xc5\x10\x7d\xfa\xbc\x39\x18\x1c\x27\xbe\x57\xca\x2d\x44\x4d\x72\x65\x3d\x1f\xc5\xb1\x8b\xbf\xab\xd1\x7f\x8b\xf2\x1e\x8d\x6a\x71\x5e\xa7\xe5\x54\xa9\x13\x76\x91\xdd\xef\xe4\xc7\x4e\x31\xd2\x4b\xf8\x79\xe5\x32\xc1\x12\xc5\xc1\xc2\x94\x95\x55\x9e\x56\xce\x95\x7b\x6c\x6b\x80\x88\xbe\xbb\x20\x1b\xc6\x59\x8d\x4a\x05\x8b\x2e\x58\x2c\x1b\xb8\x00\x53\x56\x83\xfd\xce\xd6\x49\x54\xa4\x82\x44\xff\x51\xd8\x1f\x51\x17\x05\x17\x66\x78\xd6\x46\x75\x5d\xb0\x67\x6a\x5a\x64\x6f\x94\xba\xe4\xc5\x07\xa3\xe0\xc2\x59\xa3\x93\x4b\xbc\x89\x42\x37\x31\x2b\xc9\x85\x41\x65\xd3\x8f\x17\x61\x0c\xeb\x35\x48\x81\x50\xa1\x72\x53\x27\x97\x0a\x7a\x18\x90\x16\x4c\xef\x50\xdb\x10\x5c\xa5\x4c\x4c\x3d\x4f\xef\x04\x31\x23\x96\x39\x4b\x31\x39\xef\x73\xa2\x8d\x9c\x0e\x03\x79\xdb\xc5\x10\x91\xa3\x4e\xfa\xa0\x23\xba\x38\xfa\xce\x3a\xeb\x56\x45\x0d\x4e\x25\x87\xda\x86\xf7\x13\xdc\xf0\x0c\x95\xc7\x0a\x32\x07\x4d\xfa\xb1\x4d\x81\xd6\x34\x9d\x58\xaa\x4c\xf1\x3d\x2a\x3f\xd7\xf7\x0e\x37\x30\xe3\x32\x49\x56\x76\xfa\xee\x10\x0a\xae\x8d\xf5\x05\x36\x15\x66\x1c\x45\x7a\x08\x16\xfa\x86\x9b\x74\x07\x7b\x8a\xbe\x1b\xa2\x11\x31\xb6\x8a\xa7\xae\xdb\x9a\x7f\xfc\xfd\xd5\x19\x95\xf7\xb1\xa7\x72\x29\xe5\xc8\x5c\x36\xcd\x27\xd3\x3e\x76\x0d\x6f\x14\x7d\x6a\x31\x33\xd9\x45\x76\x51\x0f\xa3\x0e\x6f\x47\xd4\xce\x79\x79\x8b\xca\xbb\x93\xe9\xa1\x6f\x13\xbd\x73\xf3\x0a\xf6\x3e\x97\xb5\x51\x34\xe7\x92\x9f\x8c\xe4\xd1\x3e\xfe\xd1\x2d\x8c\x62\x30\xd6\x75\xaa\x26\x2b\x7c\xb1\x2e\x16\x94\xdb\x8b\x23\x54\xb1\xe6\xba\xd2\xba\xdf\x5c\x5f\x69\xfb\xf8\x1b\x99\x7d\x94\xff\xa8\xe6\x9f\x92\x0f\xc9\xb1\x3f\xa6\xcc\xbd\x09\xb3\xbc\xcd\x83\x2c\xf1\x0a\x5a\xfd\x7c\x1b\x39\xd7\x0b\x02\xa7\x9a\x95\xb2\xec\x45\xd7\x0f\x91\x5d\x3f\x2c\xa7\x97\x9e\xd7\x9f\xd0\x6b\xc2\x7a\x79\xc2\xdb\xa9\xf0\x14\xdc\xf3\x42\x32\x62\x4f\x9d\xf0\x77\x2d\x45\x3f\xdc\x34\xe0\x1e\xd5\xc1\xec\x2c\xe6\xa1\x3c\xf2\x94\xd4\x99\xb9\xf9\x81\xde\x88\xba\xdc\xa0\xba\xd7\x37\x8f\x22\xe2\x49\x3c\x5b\x3f\x65\xd8\xea\x27\x8d\xdb\xf2\xd8\x46\xbf\x96\xfd\x5d\xdd\x68\xf9\xad\xba\xef\xf2\xf1\xda\x6f\x17\x2c\x06\x80\x11\x9c\x45\x15\xda\xc1\xc9\xf5\x1a\xdc\x4c\x9c\x0c\x79\x37\x2f\xdd\xda\xec\xa8\x9f\x4e\x7a\x4b\x49\x60\x6f\x3c\x69\x67\x20\xdf\x11\xeb\xad\x06\xf8\xe3\x90\xf4\xb7\xd0\xc6\xe6\x6a\xd4\x8c\x74\x19\xa0\x98\x7f\xe8\xa1\x72\x5e\xb0\xad\x57\xf1\x0a\x6f\xa1\xd1\xb7\xb2\x60\x62\x0b\x44\xe4\x31\xc6\xa0\x24\x90\x8e\x77\x41\x24\x34\x14\xcd\xe1\xb8\x38\x60\xd1\xfd\x9d\x98\x93\x12\x20\x18\x86\xca\x08\x68\x3a\xec\xfc\xf6\x6e\x1d\xdf\xa2\x31\x63\x4f\xde\xa7\xe4\x5b\x24\x20\x3f\x82\x70\x23\x1f\x2e\x1b\x2f\xf3\x23\x21\xc9\x89\xd0\x2d\x37\xbb\x7a\x93\xa4\xb2\x5c\xeb\x2a\xff\xeb\xdf\xd6\xd5\x3f\xc9\x91\x13\x1f\xdd\x21\x99\x98\x46\x71\x7f\xf8\x3c\x4a\x0d\x27\x67\xa7\xb3\x38\x7a\x06\x42\x53\x1c\x2d\xec\xbd\xac\x8b\x62\x72\x62\xd2\x46\xd5\xa9\x69\xe9\x0c\x76\x7a\x92\x3a\x3d\xfa\x2f\x7e\xb3\x07\x38\xaa\xd1\xc5\x46\x4a\x7b\x64\xb1\x36\x5c\xe2\xcd\x6d\xae\x36\xc0\xa7\x00\xb8\x99\x11\x6e\x33\xb3\x49\x7a\xd0\x6c\x61\xfa\x81\x4e\x9d\x37\x08\x7c\x2b\xa4\x42\x30\x3b\xae\x5d\x8a\xac\x80\x1b\xb8\xe1\x45\x01\xbf\xd7\xda\xc0\x06\x81\xa0\xba\x70\x07\x4b\x87\x53\x7b\x67\xf9\xf8\xfc\xbf\x60\x7e\x4e\xc1\x3f\x01\xe8\x9b\x64\x7a\x00\x6f\x12\xe7\xc5\x19\x9c\xbf\x82\x9c\x15\x1a\x27\x70\xdf\xb5\xc2\x29\xa3\xe4\xa8\x54\x4c\xfe\xeb\x99\x46\xc7\x76\x1a\x8f\x5d\xe1\x2b\xfa\xb4\xc5\x7f\x75\xeb\x99\xf3\xd2\xbd\xed\x87\xe7\xf0\x9d\x57\x74\x74\x4e\x14\xbc\xe8\xef\x3f\x6e\x1f\x5a\x58\x9a\x62\x65\xb4\x6b\x59\xf6\x90\x42\x9a\xbb\xa3\x4c\x32\x6d\x68\x13\x0f\x3d\x6a\xaf\x7d\x2a\x83\x87\x01\x31\x8d\xee\xcc\xbc\x10\x19\xdc\x77\x8d\xf1\xeb\xd5\xfb\x4b\x48\xa5\x52\x98\x9a\xe2\x00\x1a\x15\x67\x05\xff\x0f\x9d\xfb\xe6\xea\xde\x48\xa0\x1d\xbd\x99\x62\xd6\xcc\x11\xeb\xf9\x7b\x0d\x77\x01\x4e\x69\x75\x65\x8f\xe2\x21\x3d\x86\xd6\x7c\xe1\xf3\x72\x64\x3e\xa1\xc2\xc4\xf3\x8c\xc4\x34\x66\x63\xa7\xf8\x8b\x12\xcf\x78\xfe\x96\x64\x62\x70\x86\xf7\x99\x9c\x2b\x59\x4e\x8c\x9e\xad\xf9\x13\x09\xd1\xe6\xf6\xad\xc9\x9e\x29\x68\xc6\xfd\xc0\x55\xea\xab\x0b\x67\xe1\xb0\x3f\xda\xac\xe0\xfb\x66\x7a\x53\x32\x73\x51\xe2\xea\x5c\xb8\xc2\x6e\xe2\xc9\x7c\x3b\x4d\x80\xd3\x5c\x60\x22\xfb\x8a\xae\x4f\xc1\x72\x8d\x9f\xcc\xb9\xbd\x7e\x77\x77\xbf\x32\xea\x81\x0d\x9e\x62\xf7\xb4\x3d\xfe\xb1\x4a\xda\x6a\xfa\xcd\xaa\xfa\x88\xbd\x8e\xdf\x98\x86\x4b\xfc\xe1\x3b\xd3\xcc\x0d\x2d\x95\x5c\xdb\x7a\xb4\xc0\xf3\x83\xfb\x50\x04\x84\x0f\xfa\x66\xe1\x2e\x92\xbb\x6e\xe6\x12\xcd\xde\xf1\xda\x69\x24\x58\x39\xec\xf6\x5f\x0c\xe6\x48\x9d\x71\x54\x5e\xf6\xe6\x48\xe6\x50\x49\xad\xf9\xa6\xe8\xaf\x7b\xfa\xab\x26\x99\x4f\xf6\x7b\xdf\xcf\x30\x8d\x62\xf8\xf4\xf9\x08\x77\x4c\x59\x51\x21\x95\xec\x1a\xa3\xfe\xfd\x0a\x0a\x9c\xbf\xa0\x8e\xe9\x9c\x24\xab\x43\x64\xaf\x2c\x67\x29\x86\x90\x98\xb2\x3a\xfa\xdd\x7e\xd2\x9b\x71\xc9\x3b\x56\x59\x87\x40\xc9\xaa\xb1\x3f\x1d\xa0\xf0\x5f\x04\x26\x90\xc2\x07\xea\x21\x38\xbd\xcf\x82\x11\xbe\xe3\x39\xfd\x39\xf3\x2d\xe1\x1d\xab\x3e\x35\xb7\xbe\x1a\x68\xa3\xc6\xb9\x96\x97\x26\xb9\xaa\x14\x17\x26\x8f\x26\x30\x31\x7a\x9e\xc5\xe1\x0a\x9a\x38\x98\x37\xd7\x95\x8f\x35\xb8\x16\x27\x26\x27\xfd\xc7\x10\x3c\xc9\xd1\xff\x05\x00\x00\xff\xff\xa2\x96\x8a\xe1\x80\x1d\x00\x00")

func enumTmplBytes() ([]byte, error) {
	return bindataRead(
		_enumTmpl,
		"enum.tmpl",
	)
}

func enumTmpl() (*asset, error) {
	bytes, err := enumTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "enum.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xa9, 0x69, 0xdb, 0x7d, 0x72, 0xdf, 0x15, 0x91, 0x44, 0x8c, 0x7a, 0x6, 0xca, 0xcb, 0xaa, 0x28, 0xd1, 0x9a, 0xe1, 0x69, 0xe2, 0x71, 0x98, 0x47, 0xdd, 0x1f, 0xaf, 0xed, 0x6e, 0xa7, 0x8a, 0x6b}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"enum.tmpl": enumTmpl,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"enum.tmpl": {enumTmpl, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
