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

func (self *ProtonFS) OpenFile(ctx context.Context, name string, flag int, _ os.FileMode) (webdav.File, error) {
	links := self.session.Links()

	isRead := flag == os.O_RDONLY
	isWrite := flag == (os.O_RDWR | os.O_CREATE | os.O_TRUNC)

	if !isRead && !isWrite {
		return nil, webdav.ErrNotImplemented
	}

	link := links.LinkFromPath(name)
	if link == nil && isRead {
		return nil, os.ErrNotExist
	}

	if isRead {
		if link.IsDir() {
			return NewDirNode(link), nil
		}

		return nil, webdav.ErrNotImplemented
	}

	return nil, webdav.ErrNotImplemented
}

func (self *ProtonFS) RemoveAll(ctx context.Context, name string) error {
	links := self.session.Links()
	filesystem := self.session.FileSystem()

	link := links.LinkFromPath(name)
	if link == nil {
		return os.ErrNotExist
	}

	return filesystem.Delete(ctx, link)
}

func (self *ProtonFS) Rename(ctx context.Context, oldName, newName string) error {
	links := self.session.Links()
	filesystem := self.session.FileSystem()

	newName = path.Clean(newName)
	dir, file := path.Split(newName)

	link := links.LinkFromPath(oldName)
	if link == nil {
		return os.ErrNotExist
	}

	parent := links.LinkFromPath(dir)
	if parent == nil {
		return os.ErrNotExist
	}

	return filesystem.Move(ctx, link, parent, file)
}

func (self *ProtonFS) Stat(_ context.Context, name string) (os.FileInfo, error) {
	links := self.session.Links()

	link := links.LinkFromPath(name)
	if link == nil {
		return nil, os.ErrNotExist
	}

	return NewNodeInfo(link), nil
}
