package fsscanner_test

import (
	"os"
	"testing"

	fs "github.com/vhugo/fsscanner"

	"github.com/stretchr/testify/assert"
)

type mockFSFile struct {
	name    string
	content []byte
}

func mockFS(files []mockFSFile) *fs.Backend {
	f := fs.NewFS(fs.MEM)

	for _, file := range files {
		n, _ := f.Fs.OpenFile(file.name, os.O_CREATE, 0777)
		if _, err := n.Write(file.content); err != nil {
			panic(err)
		}
	}

	return f
}

func TestScanner(t *testing.T) {
	testFS := mockFS([]mockFSFile{
		{"/photo-a.jpg", []byte{0xFF, 0xD8, 0xFF}},
		{"/photo-b.jpg", []byte{0xFF, 0xD8, 0x00}},
		{"/more photos/photo-a.jpg", []byte{0xFF, 0xD8, 0xFF}},
		{"/more photos/photo-b.jpg", []byte{0xFF, 0xD8, 0x00}},
	})

	scanner, err := testFS.Scan("/")
	assert.NoError(t, err)

	for _, tc := range []struct {
		m       string
		scanner *fs.Scanner
		assert  func(*testing.T, []*fs.File, error)
	}{
		{
			m:       "no backend",
			scanner: &fs.Scanner{},
			assert: func(t *testing.T, _ []*fs.File, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			m:       "happy path",
			scanner: scanner,
			assert: func(t *testing.T, files []*fs.File, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 2, len(files))

				for _, file := range files {
					assert.Equal(t, "photo-a.jpg", file.Name)
				}
			},
		},
	} {
		t.Run(tc.m, func(t *testing.T) {
			files, err := tc.scanner.List()
			tc.assert(t, files, err)
		})
	}
}
