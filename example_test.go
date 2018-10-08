package renameonclose_test

import (
	roc "github.com/james-antill/rename-on-close"
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

	if err := nf.CloseRename(); err != nil { // Rename over the original file
		panic(err) // Dito. any failure and you keep the original.
	}
}

func ExampleCreate_diff_sync() {
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

	if d, _ := nf.IsDifferent(); d {
		// Rename over the original file, after sync'ing
		if err := nf.CloseRename(); err != nil {
			panic(err) // Note: nf.File.Sync() fails then you keep the original.
		}
	}
	// If the files are the same then the defer'd close will clean up
}
