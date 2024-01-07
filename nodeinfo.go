package main

import (
	"context"
	"github.com/StollD/proton-drive"
	"github.com/StollD/webdav"
	"io/fs"
	"os"
	"time"
)

var _ os.FileInfo = &ProtonNodeInfo{}
var _ webdav.ETager = &ProtonNodeInfo{}
var _ webdav.ContentTyper = &ProtonNodeInfo{}
var _ webdav.Hasher = &ProtonNodeInfo{}

type ProtonNodeInfo struct {
	name     string
	size     int64
	isDir    bool
	modTime  time.Time
	hash     string
	mimeType string
}

func NewNodeInfo(link *drive.Link) *ProtonNodeInfo {
	return &ProtonNodeInfo{
		name:     link.Name(),
		size:     link.Size(),
		isDir:    link.IsDir(),
		modTime:  link.ModificationTime(),
		hash:     link.ContentHash(),
		mimeType: link.MIMEType(),
	}
}

func (self *ProtonNodeInfo) Name() string {
	return self.name
}

func (self *ProtonNodeInfo) Size() int64 {
	return self.size
}

func (self *ProtonNodeInfo) Mode() fs.FileMode {
	if self.isDir {
		return 0777 | os.ModeDir
	} else {
		return 0666
	}
}

func (self *ProtonNodeInfo) ModTime() time.Time {
	return self.modTime
}

func (self *ProtonNodeInfo) IsDir() bool {
	return self.isDir
}

func (self *ProtonNodeInfo) Sys() any {
	return nil
}

func (self *ProtonNodeInfo) ETag(_ context.Context) (string, error) {
	if self.hash == "" {
		return "", webdav.ErrNotImplemented
	}

	return self.hash, nil
}

func (self *ProtonNodeInfo) ContentType(_ context.Context) (string, error) {
	if self.mimeType == "" {
		return "", webdav.ErrNotImplemented
	}

	return self.mimeType, nil
}

func (self *ProtonNodeInfo) Hashes(_ context.Context) (map[string]string, error) {
	if self.hash == "" {
		return nil, webdav.ErrNotImplemented
	}
	
	return map[string]string{
		"SHA1": self.hash,
	}, nil
}
