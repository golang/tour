// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goplay

import (
	"appengine"
	"appengine/urlfetch"
	"fmt"
	"http"
	"io"
	"os"
)

const runUrl = "http://golang.org/compile?output=json"

func init() {
	http.HandleFunc("/compile", compile)
}

func compile(w http.ResponseWriter, r *http.Request) {
	if err := passThru(w, r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Compile server error.")
	}
}

func passThru(w io.Writer, req *http.Request) os.Error {
	c := appengine.NewContext(req)
	client := urlfetch.Client(c)
	r, err := client.Post(runUrl, "text/plain", req.Body)
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
