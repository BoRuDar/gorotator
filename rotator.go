package main

import (
	"io"
	"os"
	"path/filepath"
	"sync"
)

func New(cfg Config) (io.WriteCloser, error) {
	err := checkOrCreateDir(cfg.PathToDir)
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

	return &FileRotator{
		currentFile:     f,
		currentFileSize: st.Size(),
		mu:              &sync.Mutex{},
		Config:          cfg,
	}, nil
}

type FileRotator struct {
	currentFile     *os.File
	currentFileSize int64
	mu              *sync.Mutex
	Config
}

func (r *FileRotator) Write(p []byte) (n int, err error) {
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

func (r FileRotator) Close() error {
	return r.currentFile.Close()
}

func (r *FileRotator) rotate() error {
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
	return nil
}
