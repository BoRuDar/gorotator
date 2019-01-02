# gorotator
The small lib for a file rotation.

### Quick start
```golang

package main

import (
	"fmt"
	"github.com/BoRuDa/gorotator"
)

func main() {
	fr, err := gorotator.New(gorotator.Config{
		PathToDir:        "./testdir",
		FileName:         "file.log",
		MaxFileSize:      1 * gorotator.KB,
		MamNumberOfFiles: 3,
		IsWindows:        true,
	})
	if err != nil {
		panic(err)
	}
	defer fr.Close()

	_, err = fmt.Fprintln(fr, "test")
	if err != nil {
		panic(err)
	}
}


```
