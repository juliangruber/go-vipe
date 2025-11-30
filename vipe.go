package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

const version = "0.0.0"

func main() {
	// usage
	versionPtr := flag.Bool("V", false, "print version")
	flag.Parse()
	if *versionPtr {
		fmt.Println(version)
		os.Exit(0)
	}

	// temp file
	f, err := ioutil.TempFile("", "vipe")
	check(err)

	// read from stdin
	io.Copy(f, os.Stdin)
	f.Close()

	// choose & launch editor
	editor_prog := os.getenv("VIPE_EDITOR")
	if editor_prog == nil {
		editor_prog = os.getenv("VISUAL")
	}
	if editor_prog == nil {
		editor_prog = os.Getenv("EDITOR")
	}
	if editor_prog == nil {
		editor_prog = "editor"
	}
	
	editor := exec.Command(editor_prog, f.Name())
	tty, err := os.Open("/dev/tty")
	check(err)

	editor.Stdin = tty
	editor.Stdout = os.Stdout
	editor.Stderr = os.Stderr

	check(editor.Run())

	// write to stdout
	f, err = os.Open(f.Name())
	check(err)
	io.Copy(os.Stdout, f)

	// cleanup
	f.Close()
	os.Remove(f.Name())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
