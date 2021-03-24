package fsscanner_test

import (
	"io/ioutil"
	"testing"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/stretchr/testify/assert"
)

func TestExif(t *testing.T) {

	sample, err := ioutil.ReadFile("testdata/sample.jpg")
	assert.NoError(t, err)

	testFS := mockFS([]mockFSFile{
		{"/photo-a.jpg", sample},
	})

	scan, err := testFS.Scan("/")
	assert.NoError(t, err)

	files, err := scan.List()
	assert.NoError(t, err)

	for _, file := range files {
		img, err := file.Image()
		assert.NoError(t, err)

		x, err := img.Exif()
		assert.NoError(t, err)
		assert.NotNil(t, x)

		model, errM := x.Get(exif.Model)
		assert.NoError(t, errM)

		m, e := model.StringVal()
		assert.NoError(t, e)

		assert.Equal(t, "Canon EOS 40D", m)
	}

}
