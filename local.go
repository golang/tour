// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"go/build"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"golang.org/x/tools/playground/socket"

	// Imports so that go build/install automatically installs them.
	_ "golang.org/x/tour/pic"
	_ "golang.org/x/tour/tree"
	_ "golang.org/x/tour/wc"
)

const (
	socketPath = "/socket"
)

var (
	httpListen  = flag.String("http", "127.0.0.1:3999", "host:port to listen on")
	openBrowser = flag.Bool("openbrowser", true, "open browser automatically")
)

var (
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

// findRoot is a best-effort attempt to find a tour directory
// that contains the files it needs. It may not always work.
//
// TODO: Delete after Go 1.17 is out and we can just use embed; see CL 291849.
func findRoot() (string, bool) {
	// Try finding the golang.org/x/tour package in the
	// legacy GOPATH mode workspace or in build list.
	p, err := build.Import("golang.org/x/tour", "", build.FindOnly)
	if err == nil && isRoot(p.Dir) {
		return p.Dir, true
	}
	// If that didn't work, perhaps we're not inside any module
	// and the binary was built in module mode (e.g., 'go install
	// golang.org/x/tour@latest' or 'go get golang.org/x/tour'
	// outside a module).
	// In that's the case, find out what version it is,
	// and access its content from the module cache.
	if info, ok := debug.ReadBuildInfo(); ok &&
		info.Main.Path == "golang.org/x/tour" &&
		info.Main.Replace == nil &&
		info.Main.Version != "(devel)" {
		// Make some assumptions for brevity:
		// • the 'go' binary is in $PATH
		// • the main module isn't replaced
		// • the version isn't "(devel)"
		// They should hold for the use cases we care about, until this
		// entire mechanism is obsoleted by file embedding.
		out, execError := exec.Command("go", "mod", "download", "-json", "--", "golang.org/x/tour@"+info.Main.Version).Output()
		var tourRoot struct{ Dir string }
		jsonError := json.Unmarshal(out, &tourRoot)
		if execError == nil && jsonError == nil && isRoot(tourRoot.Dir) {
			return tourRoot.Dir, true
		}
	}
	return "", false
}

func main() {
	flag.Parse()

	if os.Getenv("GAE_ENV") == "standard" {
		log.Println("running in App Engine Standard mode")
		gaeMain()
		return
	}

	// find and serve the go tour files
	root, ok := findRoot()
	if !ok {
		log.Fatalln("Couldn't find files for the Go tour. Try reinstalling it.")
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
	http.HandleFunc("/lesson/", lessonHandler)

	origin := &url.URL{Scheme: "http", Host: host + ":" + port}
	http.Handle(socketPath, socket.NewHandler(origin))

	registerStatic(root)

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

// registerStatic registers handlers to serve static content
// from the directory root.
func registerStatic(root string) {
	// Keep these static file handlers in sync with app.yaml.
	http.Handle("/favicon.ico", http.FileServer(http.Dir(filepath.Join(root, "static", "img"))))
	static := http.FileServer(http.Dir(root))
	http.Handle("/content/img/", static)
	http.Handle("/static/", static)
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

The tour server appears to be listening on an address that is
not localhost and is configured to run code snippets locally.
Anyone with access to this address and port will have access
to this machine as the user running gotour.

If you don't understand this message, hit Control-C to terminate this process.

WARNING!  WARNING!  WARNING!
`

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
var prepContent = func(r io.Reader) io.Reader { return r }

// socketAddr returns the WebSocket handler address.
var socketAddr = func() string { return "ws://" + httpAddr + socketPath }

// analyticsHTML is optional analytics HTML to insert at the beginning of <head>.
var analyticsHTML template.HTML
