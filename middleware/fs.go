package middleware

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

type FileSystem interface {
	http.FileSystem
	exist(prefix, path string) bool
}

func FS(prefix string, fs FileSystem) gin.HandlerFunc {
	fileServer := http.FileServer(fs)
	if prefix != "" {
		fileServer = http.StripPrefix(prefix, fileServer)
	}
	return func(c *gin.Context) {
		if fs.exist(prefix, c.Request.URL.Path) {
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

type LocalFileSystem struct {
	http.FileSystem
	root string
}

func LocalFile(root string) *LocalFileSystem {
	return &LocalFileSystem{
		FileSystem: gin.Dir(root, true),
		root:       root,
	}
}

func (fs *LocalFileSystem) exist(prefix, requestMapping string) bool {
	filepath := strings.TrimPrefix(requestMapping, prefix)
	filepath = path.Join(fs.root, filepath)
	stat, err := os.Stat(filepath)
	if err != nil || (stat.IsDir() && filepath != fs.root) {
		return false
	}
	return true
}

type EmbedFileSystem struct {
	http.FileSystem
}

func EmbedFile(embedFS embed.FS) *EmbedFileSystem {
	dir, _ := fs.ReadDir(embedFS, ".")
	subFS, err := fs.Sub(embedFS, dir[0].Name())
	if err != nil {
		panic(err)
	}
	return &EmbedFileSystem{
		FileSystem: http.FS(subFS),
	}
}

func (fs EmbedFileSystem) exist(prefix, path string) bool {
	_, err := fs.Open(path)
	if err != nil {
		return false
	}
	return true
}
