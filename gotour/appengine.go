// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build appengine

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"appengine"
	"appengine/urlfetch"
)

const runUrl = "http://golang.org/compile"

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/compile", compileHandler)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if err := renderTour(w, "."); err != nil {
		c.Criticalf("template render: %v", err)
	}
}

func compileHandler(w http.ResponseWriter, r *http.Request) {
	if err := passThru(w, r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Compile server error.")
	}
}

func passThru(w io.Writer, req *http.Request) error {
	c := appengine.NewContext(req)
	client := urlfetch.Client(c)
	defer req.Body.Close()
	r, err := client.Post(runUrl, req.Header.Get("Content-type"), req.Body)
	if err != nil {
		c.Errorf("making POST request:", err)
		return err
	}
	defer r.Body.Close()
	if _, err := io.Copy(w, r.Body); err != nil {
		c.Errorf("copying response Body:", err)
		return err
	}
	return nil
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
