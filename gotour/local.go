// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	// Imports so that go build/install automatically installs them.
	_ "code.google.com/p/go-tour/pic"
	_ "code.google.com/p/go-tour/tree"
	_ "code.google.com/p/go-tour/wc"
)

const basePkg = "code.google.com/p/go-tour/"

var (
	httpListen  = flag.String("http", "127.0.0.1:3999", "host:port to listen on")
	htmlOutput  = flag.Bool("html", false, "render program output as HTML")
	openBrowser = flag.Bool("openbrowser", true, "open browser automatically")
)

var (
	// a source of numbers, for naming temporary files
	uniq = make(chan int)

	// GOPATH containing the tour packages
	gopath = os.Getenv("GOPATH")
)

func isRoot(path string) bool {
	_, err := os.Stat(filepath.Join(path, "tour.article"))
	return err == nil
}

func findRoot() (string, error) {
	ctx := build.Default
	p, err := ctx.Import(basePkg, "", build.FindOnly)
	if err == nil && isRoot(p.Dir) {
		return p.Dir, nil
	}
	tourRoot := filepath.Join(runtime.GOROOT(), "misc", "tour")
	ctx.GOPATH = tourRoot
	p, err = ctx.Import(basePkg, "", build.FindOnly)
	if err == nil && isRoot(tourRoot) {
		gopath = tourRoot
		return tourRoot, nil
	}
	return "", fmt.Errorf("could not find go-tour content; check $GOROOT and $GOPATH")
}

func main() {
	flag.Parse()

	// source of unique numbers
	go func() {
		for i := 0; ; i++ {
			uniq <- i
		}
	}()

	// find and serve the go tour files
	root, err := findRoot()
	if err != nil {
		log.Fatalf("Couldn't find tour files: %v", err)
	}

	log.Println("Serving content from", root)

	fs := http.FileServer(http.Dir(root))
	http.Handle("/favicon.ico", fs)
	http.Handle("/static/", fs)
	http.Handle("/talks/", fs)

	http.HandleFunc("/compile", compileHandler)
	http.HandleFunc("/kill", killHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			if err := renderTour(w, root); err != nil {
				log.Println(err)
			}
			return
		}
		http.Error(w, "not found", 404)
	})

	host, port, err := net.SplitHostPort(*httpListen)
	if err != nil {
		log.Fatal(err)
	}
	if host == "" {
		host = "localhost"
	}
	if host != "127.0.0.1" && host != "localhost" {
		log.Print(localhostWarning)
	}

	httpAddr := host + ":" + port
	go func() {
		url := "http://" + httpAddr
		if waitServer(url) && *openBrowser && startBrowser(url) {
			log.Printf("A browser window should open. If not, please visit %s", url)
		} else {
			log.Printf("Please open your web browser and visit %s", url)
		}
	}()
	log.Fatal(http.ListenAndServe(httpAddr, nil))
}

const localhostWarning = `
WARNING!  WARNING!  WARNING!

I appear to be listening on an address that is not localhost.
Anyone with access to this address and port will have access
to this machine as the user running gotour.

If you don't understand this message, hit Control-C to terminate this process.

WARNING!  WARNING!  WARNING!
`

type response struct {
	Output string `json:"output"`
	Errors string `json:"compile_errors"`
}

func compileHandler(w http.ResponseWriter, req *http.Request) {
	resp := new(response)
	out, err := compile(req)
	if err != nil {
		if len(out) > 0 {
			resp.Errors = string(out) + "\n" + err.Error()
		} else {
			resp.Errors = err.Error()
		}
	} else {
		resp.Output = string(out)
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println(err)
	}
}

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

func killHandler(w http.ResponseWriter, r *http.Request) {
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
	cmd.Env = environ()

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

// environ returns an execution environment containing only GO* variables
// and replacing GOPATH with the value of the global var gopath.
func environ() (env []string) {
	for _, v := range os.Environ() {
		if !strings.HasPrefix(v, "GO") {
			continue
		}
		if strings.HasPrefix(v, "GOPATH=") {
			v = "GOPATH=" + gopath
		}
		env = append(env, v)
	}
	return
}

// waitServer waits some time for the http Server to start
// serving url and returns whether it starts
func waitServer(url string) bool {
	tries := 20
	for tries > 0 {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return true
		}
		time.Sleep(100 * time.Millisecond)
		tries--
	}
	return false
}

// startBrowser tries to open the URL in a browser, and returns
// whether it succeed.
func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}

// prepContent for the local tour simply returns the content as-is.
func prepContent(r io.Reader) io.Reader { return r }
