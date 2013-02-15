// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"code.google.com/p/go.talks/pkg/present"
)

func init() {
	present.PlayEnabled = true
}

// renderTour loads tour.article and the relevant HTML templates from the given
// tour root, and renders the template to the provided writer.
func renderTour(w io.Writer, root string) error {
	// Open and parse source file.
	source := filepath.Join(root, "tour.article")
	f, err := os.Open(source)
	if err != nil {
		return err
	}
	defer f.Close()
	doc, err := present.Parse(prepContent(f), source, 0)
	if err != nil {
		return err
	}

	// Set up templates.
	action := filepath.Join(root, "template", "action.tmpl")
	tour := filepath.Join(root, "template", "tour.tmpl")
	t := present.Template().Funcs(template.FuncMap{"nocode": nocode, "socketAddr": socketAddr})
	_, err = t.ParseFiles(action, tour)
	if err != nil {
		return err
	}

	// Render.
	return doc.Render(w, t)
}

// nocode returns true if the provided Section contains
// no Code elements with Play enabled.
func nocode(s present.Section) bool {
	for _, e := range s.Elem {
		if c, ok := e.(present.Code); ok && c.Play {
			return false
		}
	}
	return true
}

var commonScripts = []string{
	"jquery.js",
	"codemirror/lib/codemirror.js",
	"codemirror/lib/go.js",
	"lang.js",
}

// serveScripts registers an HTTP handler at /script.js that serves a
// concatenated set of all the scripts specified by path relative to root.
func serveScripts(root string, path ...string) error {
	modTime := time.Now()
	var buf bytes.Buffer
	scripts := append(commonScripts, path...)
	scripts = append(scripts, "tour.js")
	for _, p := range scripts {
		fn := filepath.Join(root, p)
		b, err := ioutil.ReadFile(fn)
		if err != nil {
			return err
		}
		fmt.Fprintf(&buf, "\n\n// **** %s ****\n\n", filepath.Base(fn))
		buf.Write(b)
	}
	b := buf.Bytes()
	http.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/javascript")
		http.ServeContent(w, r, "", modTime, bytes.NewReader(b))
	})
	return nil
}
