package renameonclose

import (
	"bytes"
	"fmt"
	"io"
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
	if base == "" {
		base = "."
	}
	if fname == "" {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrInvalid}
	}

	tmpfile, err := ioutil.TempFile(base, fname)
	if err != nil {
		return nil, err
	}

	return &File{File: tmpfile, oname: name}, nil
}

// Renamed returns the name of the file we will rename to
func (f *File) Renamed() string {
	return f.oname
}

// Sync the file before the rename
func (f *File) Sync() error {
	f.sync = true
	return nil
}

// CloseRename closes the file, maybe after sync'ing, and then rename's it.
func (f *File) CloseRename() error {
	f.rename = true
	return f.Close()
}

const chunckSize = 4 * 1024

// IsDifferent is the new file different to the file it will replace.
// On errors it returns true for differences.
func (f *File) IsDifferent() (bool, error) {
	nstat, err := f.Stat()
	if err != nil {
		return true, err
	}
	ostat, err := os.Stat(f.oname)
	if err != nil {
		return true, err
	}
	if nstat.Size() != ostat.Size() {
		return true, nil
	}

	of, err := os.Open(f.oname)
	if err != nil {
		return true, err
	}
	defer of.Close()

	var off int64
	b1store := make([]byte, chunckSize)
	b2store := make([]byte, chunckSize)
	for {
		n1, err1 := f.ReadAt(b1store, off)
		off += int64(n1)
		if n1 > 0 && err1 == io.EOF {
			err1 = nil
		}

		n2, err2 := of.Read(b2store)

		if err1 != nil || err2 != nil {

			if err1 == io.EOF && err2 == io.EOF {
				return false, nil
			} else if err1 == io.EOF {
				return true, nil
			} else if err2 == io.EOF {
				return true, nil
			} else {
				err := err1
				if err != nil {
					err = err2
				}
				return true, err
			}
		}

		b1 := b1store[:n1]
		b2 := b2store[:n2]
		if n1 != n2 { // Bad ?
			return false, nil
		}
		if !bytes.Equal(b1[:n1], b2[:n1]) {
			fmt.Println("JDBG b", b1, b2)

			return true, nil
		}
	}
}

// Close the file and delete
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
		return nil
	}

	return os.Remove(name)
}
