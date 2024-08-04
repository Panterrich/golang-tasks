package main

import (
	"fmt"
	"os"
)

type DirWorker interface {
	ProcessFiles(files []os.DirEntry) ([]os.DirEntry, error)
	Work(index int, file os.DirEntry) error
	Clone() DirWorker
}

type DirWalker interface {
	Walk(path string, worker DirWorker) error
}

type PreOrderDirWalker struct{}

var _ DirWalker = (*PreOrderDirWalker)(nil)

func NewPreOrderDirWalker() PreOrderDirWalker {
	return PreOrderDirWalker{}
}

func (w *PreOrderDirWalker) Walk(path string, worker DirWorker) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read dir for path %s: %v", path, err)
	}

	files, err = worker.ProcessFiles(files)
	if err != nil {
		return fmt.Errorf("failed to process files %v: %v", files, err)
	}

	for i, file := range files {
		if err := worker.Work(i, file); err != nil {
			return fmt.Errorf("failed to work file %s: %v", file.Name(), err)
		}

		if file.IsDir() {
			oldWorker := worker.Clone()

			nextPath := path + string(os.PathSeparator) + file.Name()
			if err = w.Walk(nextPath, worker); err != nil {
				return fmt.Errorf("failed to walk for path %s: %v", nextPath, err)
			}

			worker = oldWorker
		}
	}

	return nil
}
