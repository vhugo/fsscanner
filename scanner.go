package fsscanner

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/afero"
)

// Scanner represents the results for search and inpesction of media files
type Scanner struct {
	fs    afero.Fs
	root  string
	files []*File
	clock time.Time
}

// List returns all media files discovered during the scan
func (s *Scanner) List() ([]*File, error) {
	if s.fs == nil {
		return nil, fmt.Errorf("file system backend not found")
	}
	return s.files, nil
}

func newScan(fs afero.Fs, dir string) *Scanner {
	s := Scanner{
		fs:    fs,
		root:  dir,
		clock: time.Now(),
	}
	return &s
}

func (s *Scanner) walk() error {
	return afero.Walk(s.fs, s.root, s.inspect)
}

func (s *Scanner) inspect(path string, fileInfo os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	file := newFile(s.fs, path, fileInfo)

	switch {
	case file.IsImage():
		s.files = append(s.files, file)
	}

	return nil
}
