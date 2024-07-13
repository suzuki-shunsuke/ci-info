package write

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

func MkDir(fs afero.Fs, dir string) (string, error) {
	if dir == "" {
		d, err := afero.TempDir(fs, "", "ci-info")
		if err != nil {
			return "", fmt.Errorf("create a temporal directory: %w", err)
		}
		return d, nil
	}
	if !filepath.IsAbs(dir) {
		d, err := filepath.Abs(dir)
		if err != nil {
			return "", fmt.Errorf("convert -dir %s to absolute path: %w", dir, err)
		}
		dir = d
	}
	if err := fs.MkdirAll(dir, 0o755); err != nil { //nolint:mnd
		return "", fmt.Errorf("create a directory %s: %w", dir, err)
	}
	return dir, nil
}
