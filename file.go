package fsscanner

import (
	"errors"
	"os"
	"time"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
	ftypes "github.com/h2non/filetype/types"
	"github.com/spf13/afero"
)

// File represents a media file
type File struct {
	fs       afero.Fs
	Path     string
	Name     string
	ModTime  time.Time
	Size     int64
	FileInfo os.FileInfo
	Kind     types.Type
}

var (
	ErrUnsupportedType = errors.New("Unsupported type")
)

func newFile(fs afero.Fs, path string, file os.FileInfo) *File {
	f := File{
		fs:       fs,
		Path:     path,
		Name:     file.Name(),
		ModTime:  file.ModTime(),
		Size:     file.Size(),
		FileInfo: file,
	}
	return &f
}

func (f *File) Image() (*Image, error) {
	if !f.IsImage() {
		return nil, errors.New("not an image file")
	}

	img := Image{
		file: f,
	}
	return &img, nil
}

// IsImage returns whether or not a file can be considered as an image
func (f *File) IsImage() bool {
	if _, ok := matchers.Image[f.Kind]; ok {
		return true
	}

	file, _ := f.fs.Open(f.Path)
	h := make([]byte, 261)
	if _, err := file.Read(h); err != nil {
		return false
	}

	f.Kind, _ = filetype.Image(h)
	return f.Kind != ftypes.Unknown
}

func (f *File) Bytes() []byte {
	b, _ := afero.ReadFile(f.fs, f.Path)
	return b
}
