package main

import (
	"context"
	"github.com/StollD/proton-drive"
	"github.com/StollD/webdav"
	"io/fs"
	"mime"
	"path"
)

var _ webdav.File = &ProtonWriteNode{}

type ProtonWriteNode struct {
	ctx     context.Context
	session *drive.Session

	parent *drive.Link
	name   string

	writer *drive.FileWriter
}

func NewWriteNode(ctx context.Context, session *drive.Session, parent *drive.Link, name string) *ProtonWriteNode {
	return &ProtonWriteNode{
		ctx:     ctx,
		session: session,
		parent:  parent,
		name:    name,
	}
}

func (self *ProtonWriteNode) openWriter() error {
	if self.writer != nil {
		return nil
	}

	filesystem := self.session.FileSystem()

	writer, err := filesystem.Upload(self.ctx, self.parent, self.name)
	if err != nil {
		return err
	}

	self.writer = writer
	return nil
}

func (self *ProtonWriteNode) Close() error {
	if self.writer == nil {
		return nil
	}

	err := self.writer.Close()
	if err != nil {
		return err
	}

	self.writer = nil
	return nil
}

func (self *ProtonWriteNode) Read(_ []byte) (int, error) {
	return 0, webdav.ErrNotImplemented
}

func (self *ProtonWriteNode) Seek(_ int64, _ int) (int64, error) {
	return 0, webdav.ErrNotImplemented
}

func (self *ProtonWriteNode) Readdir(_ int) ([]fs.FileInfo, error) {
	return nil, webdav.ErrNotImplemented
}

func (self *ProtonWriteNode) Stat() (fs.FileInfo, error) {
	err := self.openWriter()
	if err != nil {
		return nil, err
	}

	mimeType := mime.TypeByExtension(path.Ext(self.name))
	if mimeType == "" {
		mimeType = "text/plain"
	}

	return &ProtonNodeInfo{
		name:     self.name,
		size:     self.writer.Size(),
		isDir:    false,
		modTime:  self.writer.ModTime(),
		hash:     self.writer.Hash(),
		mimeType: mimeType,
	}, nil
}

func (self *ProtonWriteNode) Write(buffer []byte) (int, error) {
	err := self.openWriter()
	if err != nil {
		return 0, err
	}

	return self.writer.Write(buffer)
}
