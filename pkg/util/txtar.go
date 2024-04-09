package util

import (
	"io/fs"
	"os"

	"github.com/holos-run/holos/pkg/errors"
	"golang.org/x/tools/txtar"
)

func MakeWalkFunc(a *txtar.Archive) fs.WalkDirFunc {
	return func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err)
		}

		if !d.IsDir() {
			if file, err := file(path); err != nil {
				return errors.Wrap(err)
			} else {
				a.Files = append(a.Files, file)
			}
		}

		return nil
	}
}

func file(path string) (file txtar.File, err error) {
	file.Name = path
	file.Data, err = os.ReadFile(path)
	return
}
