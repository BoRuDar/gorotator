package gorotator

import (
	"os"
	"regexp"
	"testing"
)

func Test_checkOrCreateDir(t *testing.T) {
	const pathToDir = "test"

	if err := checkOrCreateDir(pathToDir); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(pathToDir)
	if err != nil {
		t.Fatal("Open err: ", err)
	}

	defer f.Close()
	defer os.Remove(pathToDir)

	stat, err := f.Stat()
	if err != nil {
		t.Fatal("Stat err: ", err)
	}

	if !stat.IsDir() {
		t.Error("not a directory")
	}

	if stat.Name() != pathToDir {
		t.Error("not a directory")
	}
}

func Test_openNewFile(t *testing.T) {
	cfg := Config{
		FileName:         "test.log",
		PathToDir:        "",
		MaxFileSize:      0,
		MaxNumberOfFiles: 0,
		IsWindows:        false,
	}

	f, err := openNewFile(cfg)
	if err != nil {
		t.Fatal("openNewFile err: ", err)
	}

	defer f.Close()
	defer os.Remove(cfg.FileName)

	stat, err := f.Stat()
	if err != nil {
		t.Fatal("Stat err: ", err)
	}

	if stat.IsDir() {
		t.Error("should not be a directory")
	}

	if stat.Name() != cfg.FileName {
		t.Error("not a directory")
	}
}

func Test_newName(t *testing.T) {
	const (
		oldName        = "testFile"
		unixPattern    = `^` + oldName + `\.[0-9]{4}-[0-9]{1,2}-[0-9]{1,2}T[0-9]{1,2}:[0-9]{1,2}:[0-9]{1,2}\.[0-9]{8,9}\+[0-9]{1,2}:[0-9]{1,2}$`
		windowsPattern = `^` + oldName + `_[0-9]{4}-[0-9]{1,2}-[0-9]{1,2}T[0-9]{1,2}\.[0-9]{1,2}\.[0-9]{1,2}_[0-9]{8,9}$`
	)

	newFileName := newName(oldName, false)
	if !regexp.MustCompile(unixPattern).MatchString(newFileName) {
		t.Error(newFileName, "doesn't match the pattern:\n", unixPattern)
	}

	newFileNameWin := newName(oldName, true)
	if !regexp.MustCompile(windowsPattern).MatchString(newFileNameWin) {
		t.Error(newFileNameWin, "doesn't match the pattern:\n", windowsPattern)
	}
}
