package main

import (
	"context"
	"github.com/StollD/proton-drive"
	"github.com/StollD/webdav"
	"os"
)

var _ webdav.FileSystem = &ProtonFS{}

type ProtonFS struct {
	session *drive.Session
}

func (self *ProtonFS) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return webdav.ErrNotImplemented
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
