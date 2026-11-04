package globalshared

import (
	"fmt"
	"os"
)

func Create_dir(path string) error {
	perm := os.FileMode(0700)

	err := os.Mkdir(path, perm)
	if err != nil {
		return fmt.Errorf("failed to create root-only directory at %s: %w", path, err)
	}

	return nil
}

