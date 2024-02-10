package main

import (
	"io/fs"
	"os"

	drive "github.com/StollD/proton-drive"
	"github.com/StollD/webdav"
)

var _ webdav.File = &ProtonDirNode{}

type ProtonDirNode struct {
	info     os.FileInfo
	children []os.FileInfo
}

func NewDirNode(link *drive.Link) *ProtonDirNode {
	var children []os.FileInfo

	for child := range link.Children().Iter() {
		children = append(children, NewNodeInfo(child))
	}

	return &ProtonDirNode{
		info:     NewNodeInfo(link),
		children: children,
	}
}

func (self *ProtonDirNode) Close() error {
	return nil
}

func (self *ProtonDirNode) Read(_ []byte) (int, error) {
	return 0, webdav.ErrNotImplemented
}

func (self *ProtonDirNode) Seek(_ int64, _ int) (int64, error) {
	return 0, webdav.ErrNotImplemented
}

func (self *ProtonDirNode) Readdir(count int) ([]fs.FileInfo, error) {
	if count > 0 {
		return nil, webdav.ErrNotImplemented
	}

	return self.children, nil
}

func (self *ProtonDirNode) Stat() (fs.FileInfo, error) {
	return self.info, nil
}

func (self *ProtonDirNode) Write(_ []byte) (int, error) {
	return 0, webdav.ErrNotImplemented
}
