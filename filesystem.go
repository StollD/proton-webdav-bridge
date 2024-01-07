package main

import (
	"context"
	"github.com/StollD/proton-drive"
	"github.com/StollD/webdav"
	"os"
	"path"
)

var _ webdav.FileSystem = &ProtonFS{}

type ProtonFS struct {
	session *drive.Session
}

func (self *ProtonFS) Mkdir(ctx context.Context, name string, _ os.FileMode) error {
	links := self.session.Links()
	filesystem := self.session.FileSystem()

	name = path.Clean(name)
	dir, file := path.Split(name)

	link := links.LinkFromPath(name)
	if link != nil {
		return os.ErrExist
	}

	parent := links.LinkFromPath(dir)
	if parent == nil {
		return os.ErrNotExist
	}

	return filesystem.CreateDir(ctx, parent, file)
}

func (self *ProtonFS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	return nil, webdav.ErrNotImplemented
}

func (self *ProtonFS) RemoveAll(ctx context.Context, name string) error {
	return webdav.ErrNotImplemented
}

func (self *ProtonFS) Rename(ctx context.Context, oldName, newName string) error {
	return webdav.ErrNotImplemented
}

func (self *ProtonFS) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	return nil, webdav.ErrNotImplemented
}
