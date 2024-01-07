package main

import (
	"context"
	"github.com/StollD/proton-drive"
	"github.com/StollD/webdav"
	"io/fs"
	"os"
)

var _ webdav.File = &ProtonReadNode{}

type ProtonReadNode struct {
	ctx     context.Context
	session *drive.Session

	link *drive.Link

	info   os.FileInfo
	reader *drive.FileReader
}

func NewReadNode(ctx context.Context, session *drive.Session, link *drive.Link) *ProtonReadNode {
	return &ProtonReadNode{
		ctx:     ctx,
		session: session,
		link:    link,
		info:    NewNodeInfo(link),
	}
}

func (self *ProtonReadNode) openReader() error {
	if self.reader != nil {
		return nil
	}

	filesystem := self.session.FileSystem()

	reader, err := filesystem.Download(self.ctx, self.link)
	if err != nil {
		return err
	}

	self.reader = reader
	return nil
}

func (self *ProtonReadNode) Close() error {
	if self.reader == nil {
		return nil
	}

	err := self.reader.Close()
	if err != nil {
		return err
	}

	self.reader = nil
	return nil
}

func (self *ProtonReadNode) Read(buffer []byte) (int, error) {
	err := self.openReader()
	if err != nil {
		return 0, err
	}

	return self.reader.Read(buffer)
}

func (self *ProtonReadNode) Seek(offset int64, whence int) (int64, error) {
	err := self.openReader()
	if err != nil {
		return 0, err
	}

	return self.reader.Seek(offset, whence)
}

func (self *ProtonReadNode) Readdir(_ int) ([]fs.FileInfo, error) {
	return nil, webdav.ErrNotImplemented
}

func (self *ProtonReadNode) Stat() (fs.FileInfo, error) {
	return self.info, nil
}

func (self *ProtonReadNode) Write(_ []byte) (int, error) {
	return 0, webdav.ErrNotImplemented
}
