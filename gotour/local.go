// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	// Imports so that goinstall automatically installs them.
	_ "code.google.com/p/go-tour/pic"
	_ "code.google.com/p/go-tour/tree"
	_ "code.google.com/p/go-tour/wc"
)

const basePkg = "code.google.com/p/go-tour/"

var (
	httpListen = flag.String("http", "127.0.0.1:3999", "host:port to listen on")
	htmlOutput = flag.Bool("html", false, "render program output as HTML")
)

var (
	// a source of numbers, for naming temporary files
	uniq = make(chan int)
)

func main() {
	flag.Parse()

	// source of unique numbers
	go func() {
		for i := 0; ; i++ {
			uniq <- i
		}
	}()

	// find and serve the go tour files
	p, err := build.Default.Import(basePkg, "", build.FindOnly)
	if err != nil {
		log.Fatalf("Couldn't find tour files: %v", err)
	}
	root := p.Dir
	log.Println("Serving content from", root)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/favicon.ico" || r.URL.Path == "/" {
			fn := filepath.Join(root, "static", r.URL.Path[1:])
			http.ServeFile(w, r, fn)
			return
		}
		http.Error(w, "not found", 404)
	})
	http.Handle("/static/", http.FileServer(http.Dir(root)))
	http.Handle("/talks/", http.FileServer(http.Dir(root)))
	http.HandleFunc("/kill", kill)

	if !strings.HasPrefix(*httpListen, "127.0.0.1") &&
		!strings.HasPrefix(*httpListen, "localhost") {
		log.Print(localhostWarning)
	}

	log.Printf("Open your web browser and visit http://%s/", *httpListen)
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

var (
	commentRe = regexp.MustCompile(`(?m)^#.*\n`)
	tmpdir    string
)

func init() {
	// find real temporary directory (for rewriting filename in output)
	var err error
	tmpdir, err = filepath.EvalSymlinks(os.TempDir())
	if err != nil {
		log.Fatal(err)
	}
}

func compile(req *http.Request) (out []byte, err error) {
	stopRun()

	// x is the base name for .go, .6, executable files
	x := filepath.Join(tmpdir, "compile"+strconv.Itoa(<-uniq))
	src := x + ".go"
	bin := x
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	// rewrite filename in error output
	defer func() {
		if err != nil {
			// drop messages from the go tool like '# _/compile0'
			out = commentRe.ReplaceAll(out, nil)
		}
		out = bytes.Replace(out, []byte(src+":"), []byte("main.go:"), -1)
	}()

	// write body to x.go
	body := []byte(req.FormValue("body"))
	defer os.Remove(src)
	if err = ioutil.WriteFile(src, body, 0666); err != nil {
		return
	}

	// build x.go, creating x
	dir, file := filepath.Split(src)
	out, err = run(dir, "go", "build", "-o", bin, file)
	defer os.Remove(bin)
	if err != nil {
		return
	}

	// run x
	return run("", bin)
}

// run executes the specified command and returns its output and an error.
func run(dir string, args ...string) ([]byte, error) {
	var buf bytes.Buffer
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
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
