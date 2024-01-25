package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ReadFile(path string) ([]byte, error) {
	file, err := fileReader(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)

}

func fileReader(f string) (io.ReadCloser, error) {
	if _, err := os.Stat(f); err != nil {
		return nil, fmt.Errorf("failed to locate configuration file %q", f)
	}

	fp, err := filepath.Abs(f)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}

	return file, nil
}
