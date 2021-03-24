package fsscanner_test

import (
	"os"
	"testing"

	fs "fsscanner"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestBackend(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		for _, tc := range []struct {
			m        string
			kind     fs.Kind
			expected *fs.Backend
		}{
			{
				m:        "invalid kind",
				kind:     fs.Kind(255),
				expected: nil,
			},
			{
				m:    "afero memory map",
				kind: fs.MEM,
				expected: &fs.Backend{
					Fs: afero.NewMemMapFs(),
				},
			},
			{
				m:    "afero OS",
				kind: fs.OS,
				expected: &fs.Backend{
					Fs: afero.NewOsFs(),
				},
			},
		} {
			t.Run(tc.m, func(t *testing.T) {
				assert.Equal(t, tc.expected, fs.NewFS(tc.kind))
			})
		}
	})

	t.Run("scan", func(t *testing.T) {
		f := fs.NewFS(fs.MEM)
		assert.NotNil(t, f)

		for _, tc := range []struct {
			m        string
			dir      string
			expected error
		}{
			{
				m:        "directory not found",
				dir:      "/not_found_dir",
				expected: os.ErrNotExist,
			},
		} {
			t.Run(tc.m, func(t *testing.T) {
				_, err := f.Scan(tc.dir)
				assert.ErrorIs(t, err, tc.expected)
			})
		}
	})
}
