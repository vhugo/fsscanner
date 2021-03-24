package fsscanner

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"

	"github.com/corona10/goimagehash"
	"github.com/h2non/filetype/matchers"
	"github.com/nf/cr2"
	"github.com/rwcarlsen/goexif/exif"
	"golang.org/x/image/bmp"
)

type Image struct {
	file  *File
	exif  *exif.Exif
	phash uint64
}

// Exif returns metadata when applicable
func (i *Image) Exif() (*exif.Exif, error) {
	if i.exif != nil {
		return i.exif, nil
	}

	if _, ok := matchers.Image[i.file.Kind]; ok {
		file, err := i.file.fs.Open(i.file.Path)
		if err != nil {
			return nil, err
		}

		i.exif, err = exif.Decode(file)
		return i.exif, err
	}

	return nil, ErrUnsupportedType
}

// PerceptionHash returns a perception hash for the image in question
func (i *Image) PerceptionHash() (uint64, error) {
	r := bytes.NewReader(i.file.Bytes())

	var img image.Image
	var err error

	switch i.file.Kind.MIME.Subtype {
	case "bmp":
		img, err = bmp.Decode(r)
		if err != nil {
			return 0, err
		}

	case "jpeg":
		img, err = jpeg.Decode(r)
		if err != nil {
			return 0, err
		}

	case "tiff", "x-canon-cr2":
		img, err = cr2.Decode(r)
		if err != nil {
			return 0, err
		}

	default:
		return 0, fmt.Errorf("Unknown decode for %q", i.file.Kind.MIME.Subtype)
	}

	phash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return 0, err
	}

	i.phash = phash.GetHash()
	return i.phash, nil
}
