// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine

package main

import (
	"flag"
	"fmt"
	"go/build"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"code.google.com/p/go.tools/playground/socket"

	// Imports so that go build/install automatically installs them.
	_ "code.google.com/p/go-tour/pic"
	_ "code.google.com/p/go-tour/tree"
	_ "code.google.com/p/go-tour/wc"
)

const (
	basePkg    = "code.google.com/p/go-tour/"
	socketPath = "/socket"
)

var (
	httpListen  = flag.String("http", "127.0.0.1:3999", "host:port to listen on")
	openBrowser = flag.Bool("openbrowser", true, "open browser automatically")
)

var (
	// GOPATH containing the tour packages
	gopath = os.Getenv("GOPATH")

	httpAddr string
)

// isRoot reports whether path is the root directory of the tour tree.
// To be the root, it must have content and template subdirectories.
func isRoot(path string) bool {
	_, err := os.Stat(filepath.Join(path, "content", "welcome.article"))
	if err == nil {
		_, err = os.Stat(filepath.Join(path, "template", "index.tmpl"))
	}
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

	// find and serve the go tour files
	root, err := findRoot()
	if err != nil {
		log.Fatalf("Couldn't find tour files: %v", err)
	}

	log.Println("Serving content from", root)

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
	httpAddr = host + ":" + port

	if err := initTour(root, "SocketTransport"); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", rootHandler)
	http.Handle("/static/", http.FileServer(http.Dir(root)))
	http.HandleFunc("/lesson/", lessonHandler)
	http.Handle(socketPath, socket.Handler)

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

// rootHandler returns a handler for all the requests except the ones for lessons.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if err := renderUI(w); err != nil {
		log.Println(err)
	}
}

// lessonHandler handler the HTTP requests for lessons.
func lessonHandler(w http.ResponseWriter, r *http.Request) {
	lesson := strings.TrimPrefix(r.URL.Path, "/lesson/")
	if err := writeLesson(lesson, w); err != nil {
		if err == lessonNotFound {
			http.NotFound(w, r)
		} else {
			log.Println(err)
		}
	}
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

func init() {
	socket.Environ = environ
}

// environ returns the original execution environment with GOPATH
// replaced (or added) with the value of the global var gopath.
func environ() (env []string) {
	for _, v := range os.Environ() {
		if !strings.HasPrefix(v, "GOPATH=") {
			env = append(env, v)
		}
	}
	env = append(env, "GOPATH="+gopath)
	return
}

// waitServer waits some time for the http Server to start
// serving url. The return value reports whether it starts.
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

// socketAddr returns the WebSocket handler address.
func socketAddr() string { return "ws://" + httpAddr + socketPath }
