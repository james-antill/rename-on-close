package renameonclose_test

import (
	roc "github.com/james-antill/renameonclose"
)

func ExampleCreate() {
	content := []byte("new file's content")
	nf, err := roc.Create("abcd")
	if err != nil {
		panic(err)
	}

	defer nf.Close() // clean up

	if _, err := nf.Write(content); err != nil {
		panic(err) // This will also, try to, remove the file
	}

	if err := nf.Commit(); err != nil { // Rename over the original file
		panic(err) // Dito. any failure and you keep the original.
	}
}

func ExampleCreate_sync() {
	content := []byte("new sync file's content")
	nf, err := roc.Create("abcd")
	if err != nil {
		panic(err)
	}

	defer nf.Close() // clean up

	if _, err := nf.Write(content); err != nil {
		panic(err) // This will also, try to, remove the file
	}

	_ = nf.Sync() // This makes a sync happen at close time, so no error can happen here

	if err := nf.Commit(); err != nil { // Rename over the original file, after sync'ing
		panic(err) // Note: nf.File.Sync() fails then you keep the original.
	}
}
