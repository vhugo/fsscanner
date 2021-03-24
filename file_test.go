package fsscanner_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	sample, err := ioutil.ReadFile("testdata/sample.jpg")
	assert.NoError(t, err)

	testFS := mockFS([]mockFSFile{
		{"/photo-a.jpg", sample},
	})

	scan, err := testFS.Scan("/")
	assert.NoError(t, err)

	t.Run("is image", func(t *testing.T) {
		files, err := scan.List()
		assert.NoError(t, err)

		for _, file := range files {
			assert.True(t, file.IsImage(), "file must be an image")
		}
	})
}
