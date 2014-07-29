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

	// spawn editor
	editor := exec.Command(os.Getenv("EDITOR"), f.Name())
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
