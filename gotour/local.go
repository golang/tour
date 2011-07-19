// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"exec"
	"flag"
	"fmt"
	"go/build"
	"http"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	// Imports so that goinstall automatically installs them.
	_ "go-tour.googlecode.com/hg/pic"
	_ "go-tour.googlecode.com/hg/tree"
	_ "go-tour.googlecode.com/hg/wc"
)

const basePkg = "go-tour.googlecode.com/hg"

var (
	httpListen = flag.String("http", "127.0.0.1:3999", "host:port to listen on")
	htmlOutput = flag.Bool("html", false, "render program output as HTML")
)

var (
	// a source of numbers, for naming temporary files
	uniq = make(chan int)
	// the architecture-identifying character of the tool chain, 5, 6, or 8
	archChar string
	// where gc and ld should find the go-tour packages
	pkgDir string
)

func main() {
	flag.Parse()

	// source of unique numbers
	go func() {
		for i := 0; ; i++ {
			uniq <- i
		}
	}()

	// set archChar
	var err os.Error
	archChar, err = build.ArchChar(runtime.GOARCH)
	if err != nil {
		log.Fatal(err)
	}

	// find and serve the go tour files
	t, _, err := build.FindTree(basePkg)
	if err != nil {
		log.Fatalf("Couldn't find tour files: %v", err)
	}
	root := filepath.Join(t.SrcDir(), basePkg, "static")
	log.Println("Serving content from", root)
	//tip: http.Handle("/", http.FileServer(http.Dir(root)))
	http.Handle("/", http.FileServer(root, "/"))

	// set include path for ld and gc
	pkgDir = t.PkgDir()

	if !strings.HasPrefix(*httpListen, "127.0.0.1") &&
		!strings.HasPrefix(*httpListen, "localhost") {
		log.Print(localhostWarning)
	}

	http.HandleFunc("/kill", kill)
		
	log.Printf("Serving at http://%s/", *httpListen) 
	log.Fatal(http.ListenAndServe(*httpListen, nil))
}

const localhostWarning = `
WARNING!  WARNING!  WARNING!

I appear to be listening on an address that is not localhost.
Anyone with access to this address and port will have access
to this machine as the user running gotour.

If you don't understand this message, hit Control-C to terminate this process.

WARNING!  WARNING!  WARNING!
`

var running struct {
	sync.Mutex
	cmd *exec.Cmd
}

func stopRun() {
	running.Lock()
	if running.cmd != nil {
		running.cmd.Process.Kill()
		running.cmd = nil
	}
	running.Unlock()
}

func kill(w http.ResponseWriter, r *http.Request) {
	stopRun()
}

func compile(req *http.Request) (out []byte, err os.Error) {
	stopRun()

	// x is the base name for .go, .6, executable files
	x := os.TempDir() + "/compile" + strconv.Itoa(<-uniq)
	src := x + ".go"
	obj := x + "." + archChar
	bin := x
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	// rewrite filename in error output
	defer func() {
		out = bytes.Replace(out, []byte(src+":"), []byte("main.go:"), -1)
	}()

	// write body to x.go
	body := new(bytes.Buffer)
	if _, err = body.ReadFrom(req.Body); err != nil {
		return
	}
	if err = ioutil.WriteFile(src, body.Bytes(), 0666); err != nil {
		return
	}

	// build x.go, creating x.6
	out, err = run(archChar+"g", "-I", pkgDir, "-o", obj, src)
	defer os.Remove(obj)
	if err != nil {
		return
	}

	// link x.6, creating x (the program binary)
	out, err = run(archChar+"l", "-L", pkgDir, "-o", bin, obj)
	defer os.Remove(bin)
	if err != nil {
		return
	}

	// run x
	return run(bin)
}

// run executes the specified command and returns its output and an error.
func run(args ...string) ([]byte, os.Error) {
	var buf bytes.Buffer
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &buf
	cmd.Stderr = cmd.Stdout
	
	// Start command and leave in 'running'.
	running.Lock()
	if running.cmd != nil {
		defer running.Unlock()
		return nil, fmt.Errorf("already running %s", running.cmd.Path)
	}
	if err := cmd.Start(); err != nil {
		running.Unlock()
		return nil, err
	}
	running.cmd = cmd
	running.Unlock()
	
	// Wait for the command.  Clean up,
	err := cmd.Wait()
	running.Lock()
	if running.cmd == cmd {
		running.cmd = nil
	}
	running.Unlock()
	return buf.Bytes(), err
}
