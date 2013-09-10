// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
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

var (
	article     = flag.String("article", "tour.article", "article to load for the tour")
	tourContent []byte
)

// initTour loads tour.article and the relevant HTML templates from the given
// tour root, and renders the template to the tourContent global variable.
func initTour(root string) error {
	// Make sure playground is enabled before rendering.
	present.PlayEnabled = true

	// Open and parse source file.
	source := *article
	f, err := os.Open(source)
	if err != nil {
		// See if it exists in the root.
		source = filepath.Join(root, "tour.article")
		f, err = os.Open(source)
		if err != nil {
			return err
		}
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
	buf := new(bytes.Buffer)
	if err := doc.Render(buf, t); err != nil {
		return err
	}
	tourContent = buf.Bytes()
	return nil
}

// renderTour writes the tour content to the provided Writer.
func renderTour(w io.Writer) error {
	if tourContent == nil {
		panic("renderTour called before successful initTour")
	}
	_, err := w.Write(tourContent)
	return err
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

var scripts = []string{
	"jquery.js",
	"codemirror/lib/codemirror.js",
	"codemirror/lib/go.js",
	"lang.js",
	"playground.js",
	"tour.js",
}

// serveScripts registers an HTTP handler at /scripts.js that serves all the
// scripts specified by the variable above, and appends a line that initializes
// the tour with the specified transport.
func serveScripts(root, transport string) error {
	modTime := time.Now()
	var buf bytes.Buffer
	for _, p := range scripts {
		fn := filepath.Join(root, p)
		b, err := ioutil.ReadFile(fn)
		if err != nil {
			return err
		}
		fmt.Fprintf(&buf, "\n\n// **** %s ****\n\n", filepath.Base(fn))
		buf.Write(b)
	}
	fmt.Fprintf(&buf, "\ninitTour(new %v());\n", transport)
	b := buf.Bytes()
	http.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/javascript")
		http.ServeContent(w, r, "", modTime, bytes.NewReader(b))
	})
	return nil
}
