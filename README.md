# gorotator
The small lib for a file rotation.

### Quick start
```golang

  fr, err := New(Config{
		PathToDir:        "./testdir",
		FileName:         "file.log",
		MaxFileSize:      1 * BytesInKB,
		MamNumberOfFiles: 3,
		IsWindows:        true,
	})
  defer fr.Close()

```
