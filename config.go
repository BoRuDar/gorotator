package gorotator

type Config struct {
	FileName         string
	PathToDir        string
	MaxFileSize      int64 // in bytes
	MamNumberOfFiles int
	IsWindows        bool
}
