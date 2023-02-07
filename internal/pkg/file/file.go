package file

import "io"

type File struct {
	r         io.Reader
	Name      string
	Extension string
	Size      int
}

func New(r io.Reader, name string, extension string, size int) *File {
	return &File{
		r:         r,
		Name:      name,
		Extension: extension,
		Size:      size,
	}
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.r.Read(p)
}
