// Package utils common utils, file util
package utils

import (
	"io"
	"os"
	"time"

	"github.com/xwi88/kit4go/json"
)

// IsExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// IsDir checks whether the path is directory or not.
func IsDir(path string) bool {
	s, err := os.Stat(path)
	return err == nil && s.IsDir()
}

// IsDir checks whether the path is file or not.
func IsFile(path string) bool {
	s, err := os.Stat(path)
	return err == nil && !s.IsDir()
}

// CopyFile copy file from src to dst
func CopyFile(dstFile string, srcFile string) (written int64, err error) {
	src, err := os.Open(srcFile)
	if err != nil {
		return 0, err
	}
	defer src.Close()

	dst, err := os.OpenFile(dstFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

// GetFileInfo get the file info, if not exist, return ""
func GetFileInfo(file string) string {
	s, err := os.Stat(file)
	if err != nil {
		return ""
	}
	type fileInfo struct {
		Name    string
		Size    int64
		Mode    uint32
		ModeStr string
		ModTime time.Time
		IsDir   bool
	}
	info := fileInfo{
		Name:    s.Name(),
		Size:    s.Size(),
		Mode:    uint32(s.Mode()),
		ModeStr: s.Mode().String(),
		ModTime: s.ModTime(),
		IsDir:   s.IsDir(),
	}
	bts, _ := json.Marshal(info)
	return string(bts)
}
