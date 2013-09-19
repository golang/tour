// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build appengine

package main

import (
	"bufio"
	"bytes"
	"io"
	"net/http"

	"appengine"

	_ "code.google.com/p/go.tools/playground"
)

const runUrl = "http://golang.org/compile"

func init() {
	http.HandleFunc("/", rootHandler)
	err := serveScripts("js", "HTTPTransport")
	if err != nil {
		panic(err)
	}
	if err := initTour("."); err != nil {
		panic(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if err := renderTour(w); err != nil {
		c.Criticalf("template render: %v", err)
	}
}

// prepContent returns a Reader that produces the content from the given
// Reader, but strips the prefix "#appengine: " from each line. It also drops
// any non-blank like that follows a series of 1 or more lines with the prefix.
func prepContent(in io.Reader) io.Reader {
	var prefix = []byte("#appengine: ")
	out, w := io.Pipe()
	go func() {
		r := bufio.NewReader(in)
		drop := false
		for {
			b, err := r.ReadBytes('\n')
			if err != nil && err != io.EOF {
				w.CloseWithError(err)
				return
			}
			if bytes.HasPrefix(b, prefix) {
				b = b[len(prefix):]
				drop = true
			} else if drop {
				if len(b) > 1 {
					b = nil
				}
				drop = false
			}
			if len(b) > 0 {
				w.Write(b)
			}
			if err == io.EOF {
				w.Close()
				return
			}
		}
	}()
	return out
}

// socketAddr returns the WebSocket handler address.
// The App Engine version does not provide a WebSocket handler.
func socketAddr() string { return "" }
