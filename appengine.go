// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	_ "golang.org/x/tools/playground"
)

func gaeMain() {
	prepContent = gaePrepContent
	socketAddr = gaeSocketAddr
	analyticsHTML = template.HTML(os.Getenv("TOUR_ANALYTICS"))

	if err := initTour(".", "HTTPTransport"); err != nil {
		log.Fatal(err)
	}

	http.Handle("/", hstsHandler(rootHandler))
	http.Handle("/lesson/", hstsHandler(lessonHandler))

	registerStatic(".")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// gaePrepContent returns a Reader that produces the content from the given
// Reader, but strips the prefix "#appengine:", optionally followed by a space, from each line.
// It also drops any non-blank line that follows a series of 1 or more lines with the prefix.
func gaePrepContent(in io.Reader) io.Reader {
	var prefix = []byte("#appengine:")
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
				if b[0] == ' ' {
					// Consume a single space after the prefix.
					b = b[1:]
				}
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

// gaeSocketAddr returns the WebSocket handler address.
// The App Engine version does not provide a WebSocket handler.
func gaeSocketAddr() string { return "" }

// hstsHandler wraps an http.HandlerFunc such that it sets the HSTS header.
func hstsHandler(fn http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; preload")
		fn(w, r)
	})
}
