package web

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var DistFS embed.FS

// GetDistFS 获取嵌入的前端静态文件系统
func GetDistFS() (fs.FS, error) {
	return fs.Sub(DistFS, "dist")
}

// GetDistFSRoot 获取完整的嵌入文件系统（包含dist目录）
func GetDistFSRoot() fs.FS {
	return DistFS
}
