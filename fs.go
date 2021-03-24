package fsscanner

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// Backend represents the adapter to access data
type Backend struct {
	Fs afero.Fs
}

// Kind represents the type of backend to choose from when using a new
type Kind uint8

const (
	MEM Kind = iota
	OS
)

// New returns a fresh instance for a specific backend type
func NewFS(k Kind) *Backend {
	var af afero.Fs

	switch k {
	case MEM:
		af = afero.NewMemMapFs()

	case OS:
		af = afero.NewOsFs()

	default:
		return nil
	}

	f := Backend{
		Fs: af,
	}
	return &f
}

// Scan searches media files recursivily from one directory
func (b *Backend) Scan(dir string) (*Scanner, error) {
	s := newScan(b.Fs, dir)
	return s, s.walk()
}

// AvailableFileName checks whether a filename name is used and append a number
// to the proposed file name
func (b *Backend) AvailableFileName(dir, filename string) string {
	return b.availableFileName(0, dir, filename)
}

func (b *Backend) availableFileName(v int, dir, filename string) string {
	var newfilename string
	newfilename = filename

	if v > 0 {
		extension := filepath.Ext(filename)
		newfilename = strings.Replace(filename, extension, fmt.Sprintf(".%04d%s", v, extension), 1)
	}

	fpath := filepath.Join(dir, newfilename)
	if _, err := b.Fs.Stat(fpath); err == nil {
		v++
		// fmt.Println(v)
		return b.availableFileName(v, dir, filename)
	}
	return fpath
}
