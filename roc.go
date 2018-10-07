package renameonclose

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// File is a wrapper around os.File, but can sync and then rename/remove on close.
type File struct {
	*os.File
	oname  string
	sync   bool
	rename bool
	closed bool
}

// Create a new temporary file (via. ioutil.TempFile) that on close
// will either can rename over the path given or delete itself.
func Create(name string) (f *File, err error) {

	base, fname := filepath.Split(name)
	tmpfile, err := ioutil.TempFile(base, fname)
	if err != nil {
		return nil, err
	}

	return &File{File: tmpfile, oname: name}, nil
}

// Renamed the name of the file we will rename to
func (f *File) Renamed() string {
	return f.oname
}

// Sync the file before the rename
func (f *File) Sync() error {
	f.sync = true
	return nil
}

// Commit the file via. a sync/close/rename.
func (f *File) Commit() error {
	f.rename = true
	return f.Close()
}

// Close the file, maybe after sync'ing, and then either rename or delete
func (f *File) Close() error {

	if f == nil {
		return os.ErrInvalid
	}
	name := f.Name()

	if f.closed {
		return &os.PathError{Op: "close", Path: name, Err: os.ErrClosed}
	}
	f.closed = true

	if f.sync && f.rename {
		if err := f.Sync(); err != nil {
			os.Remove(name)
			return err
		}
	}

	if err := f.File.Close(); err != nil {
		os.Remove(name)
		return err
	}

	if f.rename {
		if err := os.Rename(name, f.oname); err != nil {
			os.Remove(name)
			return err
		}
	}

	return os.Remove(name)
}
