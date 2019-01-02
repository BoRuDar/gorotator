package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
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
		return fmt.Sprintf("%s.%s", oldName, t.Format(time.RFC3339Nano))
	}

	// Windows doesn't support `:` in file names
	// file will look like `file.log_2019-01-02T21.38.26_369788600`
	return fmt.Sprintf(
		"%s_%sT%d.%d.%d_%d",
		oldName, t.Format("2006-01-02"), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(),
	)
}

func getRotatableFiles(cfg Config) (sliceOfFiles []string, err error) {
	absPathToDir, err := filepath.Abs(cfg.PathToDir)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(absPathToDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		fileName := info.Name()
		if strings.Contains(fileName, cfg.FileName) && fileName != cfg.FileName {
			sliceOfFiles = append(sliceOfFiles, fileName)
		}

		return nil
	})

	return
}
