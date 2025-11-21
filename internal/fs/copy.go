package fs

import (
	"io"
	"os"
	"path/filepath"
)

func CopyCodebase(src, dest string) ([]string, error) {
	var copied []string

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dest, rel)

		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		// Copy file
		if err := copyFile(path, target); err != nil {
			return err
		}

		copied = append(copied, target)
		return nil
	})

	return copied, err
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
