package util

import (
	"github.com/holos-run/holos/pkg/wrapper"
	"golang.org/x/tools/txtar"
	"io/fs"
	"os"
)

func MakeWalkFunc(a *txtar.Archive) fs.WalkDirFunc {
	return func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return wrapper.Wrap(err)
		}

		if !d.IsDir() {
			if file, err := file(path); err != nil {
				return wrapper.Wrap(err)
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
