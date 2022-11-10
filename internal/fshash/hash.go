package fshash

import (
	"fmt"
	"hash"
	"io"
	"io/fs"
)

// Of writes filesystem contents to the hasher.
func Of(hasher hash.Hash, fsys fs.FS) error {
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if _, err := io.WriteString(hasher, path); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if !d.Type().IsRegular() {
			return nil
		}
		if err := hashFile(hasher, fsys, path); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		return nil
	})
}

func hashFile(hasher hash.Hash, fsys fs.FS, path string) error {
	var f, errOpen = fsys.Open(path)
	if errOpen != nil {
		return errOpen
	}
	defer f.Close()

	var _, errHash = io.Copy(hasher, f)
	return errHash
}
