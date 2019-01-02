package main

import (
	"fmt"
	"time"
)

func main() {
	fr, err := New(Config{
		PathToDir:        "./testdir",
		FileName:         "file.log",
		MaxFileSize:      2 * BytesInKB,
		MamNumberOfFiles: 3,
		IsWindows:        true,
	})
	defer fr.Close()

	for i := 0; i < 2000; i++ {
		_, err = fr.Write([]byte(fmt.Sprintf("test line #%d\n", i)))
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Microsecond * 100)
	}
}
