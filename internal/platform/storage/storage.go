package storage

import (
	"context"
	"io"
)

type File struct {
	Content     io.Reader
	Filename    string
	ContentType string
	Size        int64
}

type Storage interface {
	Upload(ctx context.Context, key string, file *File) error
	Delete(ctx context.Context, key string) error
}
