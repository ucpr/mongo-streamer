package persistent

import (
	"context"
	"io"
	"os"
)

type File struct {
	file *os.File
}

var _ Storage = (*File)(nil)

func NewFileWriter(filepath string) (*File, error) {
	// TODO: check dir exists
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, err
	}

	return &File{file: f}, nil
}

func (w *File) Write(s string) error {
	_, err := w.file.WriteString(s)
	return err
}

func (w *File) Clear() error {
	// clear data by overwriting with an empty string
	_, err := w.file.WriteString("")
	return err
}

func (w *File) Read() (string, error) {
	b, err := io.ReadAll(w.file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (w *File) Close(ctx context.Context) error {
	return w.file.Close()
}
