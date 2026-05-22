package upload

import "io"

type File struct {
	Content     io.Reader
	Filename    string
	ContentType string
	Size        int64
}
