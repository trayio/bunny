package temporary

import (
	"os"
	"testing"
)

const destination = "/tmp/letshopethisdoesntexist"

func TestNewFile(t *testing.T) {
	f, err := NewFile()
	if err != nil {
		t.Error("Failed to create temporary file:", err)
	}

	f.File.Close()
	os.Remove(f.File.Name())
}

func TestMove(t *testing.T) {
	f, err := NewFile()
	if err != nil {
		t.Error("Failed to create temporary file:", err)
	}

	if err := f.Move(destination); err != nil {
		t.Error("Failed to move file:", err)
	}

	if _, err := os.Stat(destination); os.IsNotExist(err) {
		t.Errorf(`Destination "%s" doesn't exist: %s\n`, err)
	}
}
