package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	roc "github.com/james-antill/rename-on-close"
)

var syncFlag bool
var diffFlag bool

func init() {
	flag.BoolVar(&syncFlag, "sync", false, "Sync the file before close")
	flag.BoolVar(&diffFlag, "diff", false, "Diff the file and only mv over if it's changed")
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		panic("Format: <ofile> <nfile>")
	}
	nf, err := roc.Create(args[1])
	if err != err {
		panic(err)
	}
	defer nf.Close()

	of, err := os.Open(args[0])
	if err != nil {
		panic(err)
	}
	defer of.Close()

	if syncFlag {
		nf.Sync()
	}
	if _, err := io.Copy(nf, of); err != nil {
		panic(err)
	}
	if diffFlag {
		if d, _ := nf.IsDifferent(); d {
			fmt.Println("Using new file:", nf.Name(), "for:", nf.Renamed())
			err = nf.CloseRename()
		} else {
			fmt.Println("Using original file:", nf.Renamed())
		}
	} else {
		fmt.Println("Using new file:", nf.Name(), "for:", nf.Renamed())
		err = nf.CloseRename()
	}
	if err != nil {
		panic(err)
	}
}
