package gorotator

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

func New(cfg Config) (io.WriteCloser, error) {
	err := checkOrCreateDir(cfg.PathToDir)
	if err != nil {
		return nil, err
	}

	sliceOfFiles, err := getRotatableFiles(cfg)
	if err != nil {
		return nil, err
	}

	f, err := openNewFile(cfg)
	if err != nil {
		return nil, err
	}

	st, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return &fileRotator{
		currentFile:     f,
		currentFileSize: st.Size(),
		mu:              &sync.Mutex{},
		filesToWatch:    sliceOfFiles,
		Config:          cfg,
	}, nil
}

type fileRotator struct {
	currentFileSize int64
	currentFile     *os.File
	filesToWatch    []string
	mu              *sync.Mutex
	Config
}

func (r *fileRotator) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.currentFileSize+int64(len(p)) >= r.MaxFileSize {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = r.currentFile.Write(p)
	r.currentFileSize += int64(n)
	return
}

func (r fileRotator) Close() error {
	return r.currentFile.Close()
}

func (r *fileRotator) rotate() error {
	err := r.currentFile.Close()
	if err != nil {
		return err
	}

	oldFilePath, err := filepath.Abs(r.currentFile.Name())
	if err != nil {
		return err
	}

	newFileName := newName(r.FileName, r.IsWindows)
	newFilePath, err := filepath.Abs(filepath.Join(r.PathToDir, newFileName))
	if err != nil {
		return err
	}

	r.filesToWatch = append(r.filesToWatch, newFileName)

	err = os.Rename(oldFilePath, newFilePath)
	if err != nil {
		return err
	}

	r.currentFile, err = openNewFile(r.Config)
	if err != nil {
		return err
	}

	s, err := r.currentFile.Stat()
	if err != nil {
		return err
	}
	r.currentFileSize = s.Size()

	err = r.checkOrDelete()
	if err != nil {
		return err
	}

	return nil
}

func (r *fileRotator) checkOrDelete() error {
	if len(r.filesToWatch)+1 <= r.MaxNumberOfFiles {
		return nil // nothing to do
	}

	sort.Slice(r.filesToWatch, func(i, j int) bool {
		return r.filesToWatch[i] < r.filesToWatch[j]
	})

	absPathToFile, err := filepath.Abs(filepath.Join(r.PathToDir, r.filesToWatch[0]))
	if err != nil {
		return err
	}

	err = os.Remove(absPathToFile)
	if err != nil {
		return err
	}

	r.filesToWatch = r.filesToWatch[1:]
	return nil
}
