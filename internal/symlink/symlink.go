package symlink

import (
	"os"
	"path/filepath"
)

// ReadLink reads the target of a symbolic link
func ReadLink(symlinkPath string) (string, error) {
	target, err := os.Readlink(symlinkPath)
	if err != nil {
		return "", err
	}

	// Convert to absolute path if relative
	if !filepath.IsAbs(target) {
		dir := filepath.Dir(symlinkPath)
		target = filepath.Join(dir, target)
	}

	return target, nil
}

// Create creates a symbolic link from target to symlinkPath
func Create(target, symlinkPath string) error {
	// Remove existing symlink if it exists
	if _, err := os.Lstat(symlinkPath); err == nil {
		if err := os.Remove(symlinkPath); err != nil {
			return err
		}
	}

	// Create the symbolic link
	return os.Symlink(target, symlinkPath)
}
