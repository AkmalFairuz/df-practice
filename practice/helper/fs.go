package helper

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// RemoveDir removes a directory and its contents recursively.
func RemoveDir(dir string) error {
	return os.RemoveAll(dir)
}

// CopyDir copies a directory and its contents recursively.
func CopyDir(src, dst string) error {
	// Get properties of the source directory.
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("could not get source directory info: %v", err)
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	// Create the destination directory.
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("could not create destination directory: %v", err)
	}

	// Read contents of the source directory.
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("could not read source directory: %v", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory.
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file.
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile copies a single file from src to dst.
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("could not create destination file: %v", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}

	// Copy file permissions.
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("could not get source file info: %v", err)
	}
	return os.Chmod(dst, srcInfo.Mode())
}
