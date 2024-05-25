package util

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type MyLocalFileSystem struct {
	http.FileSystem
	fs           static.ServeFileSystem
	excludePaths []string
}

func MyLocalFile(root string, indexes bool, epaths []string) *MyLocalFileSystem {
	return &MyLocalFileSystem{
		FileSystem:   gin.Dir(root, indexes),
		fs:           static.LocalFile(root, indexes),
		excludePaths: epaths,
	}
}

func (p *MyLocalFileSystem) Exists(prefix string, filepath string) bool {
	// fmt.Printf("prefix:%v filepath:%v\n", prefix, filepath)
	for _, excludePrefix := range p.excludePaths {
		if p := strings.TrimPrefix(filepath, excludePrefix); len(p) < len(filepath) {
			fmt.Printf("filePath:%v hit exclude prefix:%v", filepath, excludePrefix)
			return false
		}
	}
	return p.fs.Exists(prefix, filepath)
}
