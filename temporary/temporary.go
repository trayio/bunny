package temporary

import (
	"io/ioutil"
	"os"
)

type temporaryFile struct {
	File *os.File
}

// NewFile creates a new temporary file and return temporaryFile struct
func NewFile() (*temporaryFile, error) {
	f, err := ioutil.TempFile("", "bunny")
	if err != nil {
		return nil, err
	}
	return &temporaryFile{File: f}, nil
}

// Move moves current temporary file to a new destination.
func (t *temporaryFile) Move(destination string) error {
	if err := t.File.Sync(); err != nil {
		return err
	}

	if err := t.File.Close(); err != nil {
		return err
	}

	if err := os.Rename(t.File.Name(), destination); err != nil {
		return err
	}

	if err := os.Chmod(destination, 0444); err != nil {
		return err
	}

	return nil
}
