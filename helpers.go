package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"
)

func checkOrCreateDir(pathToDir string) (err error) {
	absPathToDir, err := filepath.Abs(pathToDir)
	if err != nil {
		return err
	}

	dir, err := os.OpenFile(absPathToDir, os.O_RDONLY, 0666)
	if err != nil {
		return os.MkdirAll(absPathToDir, 0666)
	}
	return dir.Close()
}

func openNewFile(cfg Config) (f *os.File, err error) {
	absPathToFile, err := filepath.Abs(path.Join(cfg.PathToDir, cfg.FileName))
	if err != nil {
		return nil, err
	}
	return os.OpenFile(absPathToFile, os.O_APPEND|os.O_CREATE, 0666)
}

func newName(oldName string, isWindows bool) string {
	t := time.Now()

	if !isWindows {
		return oldName + "." + time.Now().Format(time.RFC3339)
	}
	return fmt.Sprintf("%s_%sT%d.%d.%d", oldName, t.Format("2006-01-02"), t.Hour(), t.Minute(), t.Second())
}
