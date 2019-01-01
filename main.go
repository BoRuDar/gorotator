package main

import (
	"fmt"
)

func main() {
	fr, err := New(Config{
		PathToDir:        "./testdir",
		FileName:         "file.log",
		MaxFileSize:      1 * BytesInKB,
		MamNumberOfFiles: 3,
		IsWindows:        true,
	})
	defer fr.Close()

	_, err = fmt.Fprintf(fr, "%+v - %v", fr, err)

	for i := 0; i < 100; i++ {
		_, err = fr.Write([]byte(fmt.Sprintf("test line #%d\n", i)))
		if err != nil {
			panic(err)
		}
	}
}
