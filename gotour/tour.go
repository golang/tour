// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"code.google.com/p/go.tools/godoc/static"
	"code.google.com/p/go.tools/present"
)

var (
	uiContent      []byte
	lessons        = make(map[string][]byte)
	lessonNotFound = fmt.Errorf("lesson not found")
)

// initTour loads tour.article and the relevant HTML templates from the given
// tour root, and renders the template to the tourContent global variable.
func initTour(root, transport string) error {
	// Make sure playground is enabled before rendering.
	present.PlayEnabled = true

	// Set up templates.
	action := filepath.Join(root, "template", "action.tmpl")
	tmpl, err := present.Template().ParseFiles(action)
	if err != nil {
		return fmt.Errorf("parse templates: %v", err)
	}

	// Init lessons.
	contentPath := filepath.Join(root, "content")
	if err := initLessons(tmpl, contentPath); err != nil {
		return fmt.Errorf("init lessons: %v", err)
	}

	// Init UI
	index := filepath.Join(root, "template", "index.tmpl")
	ui, err := template.ParseFiles(index)
	if err != nil {
		return fmt.Errorf("parse index.tmpl: %v", err)
	}
	buf := new(bytes.Buffer)

	data := struct {
		SocketAddr string
		Transport  template.JS
	}{socketAddr(), template.JS(transport)}

	if err := ui.Execute(buf, data); err != nil {
		return fmt.Errorf("render UI: %v", err)
	}
	uiContent = buf.Bytes()

	return initScript(root)
}

// initLessonss finds all the lessons in the passed directory, renders them,
// using the given template and saves the content in the lessons map.
func initLessons(tmpl *template.Template, content string) error {
	dir, err := os.Open(content)
	if err != nil {
		return err
	}
	files, err := dir.Readdirnames(0)
	if err != nil {
		return err
	}
	for _, f := range files {
		if filepath.Ext(f) != ".article" {
			continue
		}
		content, err := parseLesson(tmpl, filepath.Join(content, f))
		if err != nil {
			return fmt.Errorf("parsing %v: %v", f, err)
		}
		name := strings.TrimSuffix(f, ".article")
		lessons[name] = content
	}
	return nil
}

// File defines the JSON form of a code file in a page.
type File struct {
	Name    string
	Content string
}

// Page defines the JSON form of a tour lesson page.
type Page struct {
	Title   string
	Content string
	Files   []File
}

// Lesson defines the JSON form of a tour lesson.
type Lesson struct {
	Title       string
	Description string
	Pages       []Page
}

// parseLesson parses and returns a lesson content given its name and
// the template to render it.
func parseLesson(tmpl *template.Template, path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	doc, err := present.Parse(prepContent(f), path, 0)
	if err != nil {
		return nil, err
	}

	lesson := Lesson{
		doc.Title,
		doc.Subtitle,
		make([]Page, len(doc.Sections)),
	}

	for i, sec := range doc.Sections {
		p := &lesson.Pages[i]
		w := new(bytes.Buffer)
		if err := sec.Render(w, tmpl); err != nil {
			return nil, fmt.Errorf("render section: %v", err)
		}
		p.Title = sec.Title
		p.Content = w.String()
		codes := findPlayCode(sec)
		p.Files = make([]File, len(codes))
		for i, c := range codes {
			f := &p.Files[i]
			f.Name = c.FileName
			f.Content = string(c.Raw)
		}
	}

	w := new(bytes.Buffer)
	if err := json.NewEncoder(w).Encode(lesson); err != nil {
		return nil, fmt.Errorf("encode lesson: %v", err)
	}
	return w.Bytes(), nil
}

// findPlayCode returns a slide with all the Code elements in the given
// Elem with Play set to true.
func findPlayCode(e present.Elem) []*present.Code {
	var r []*present.Code
	switch v := e.(type) {
	case present.Code:
		if v.Play {
			r = append(r, &v)
		}
	case present.Section:
		for _, s := range v.Elem {
			r = append(r, findPlayCode(s)...)
		}
	}
	return r
}

// writeLesson writes the tour content to the provided Writer.
func writeLesson(name string, w io.Writer) error {
	if uiContent == nil {
		panic("writeLesson called before successful initTour")
	}
	if len(name) == 0 {
		return writeAllLessons(w)
	}
	l, ok := lessons[name]
	if !ok {
		return lessonNotFound
	}
	_, err := w.Write(l)
	return err
}

func writeAllLessons(w io.Writer) error {
	if _, err := fmt.Fprint(w, "{"); err != nil {
		return err
	}
	nLessons := len(lessons)
	for k, v := range lessons {
		if _, err := fmt.Fprintf(w, "%q:%s", k, v); err != nil {
			return err
		}
		nLessons--
		if nLessons != 0 {
			if _, err := fmt.Fprint(w, ","); err != nil {
				return err
			}
		}
	}
	_, err := fmt.Fprint(w, "}")
	return err
}

// renderUI writes the tour UI to the provided Writer.
func renderUI(w io.Writer) error {
	if uiContent == nil {
		panic("renderUI called before successful initTour")
	}
	_, err := w.Write(uiContent)
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

// initScript concatenates all the javascript files needed to render
// the tour UI and serves the result on /script.js.
func initScript(root string) error {
	modTime := time.Now()
	b := new(bytes.Buffer)

	content, ok := static.Files["playground.js"]
	if !ok {
		return fmt.Errorf("playground.js not found in static files")
	}
	b.WriteString(content)

	// Keep this list in dependency order
	files := []string{
		"static/lib/jquery.min.js",
		"static/lib/jquery-ui.min.js",
		"static/lib/angular.min.js",
		"static/lib/codemirror/lib/codemirror.js",
		"static/lib/codemirror/mode/go/go.js",
		"static/lib/angular-ui.min.js",
		"static/js/app.js",
		"static/js/controllers.js",
		"static/js/directives.js",
		"static/js/services.js",
		"static/js/values.js",
	}

	for _, file := range files {
		f, err := ioutil.ReadFile(filepath.Join(root, file))
		if err != nil {
			return fmt.Errorf("couldn't open %v: %v", file, err)
		}
		_, err = b.Write(f)
		if err != nil {
			return fmt.Errorf("error concatenating %v: %v", file, err)
		}
	}

	http.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/javascript")
		// Set expiration time in one week.
		w.Header().Set("Cache-control", "max-age=604800")
		http.ServeContent(w, r, "", modTime, bytes.NewReader(b.Bytes()))
	})

	return nil
}
