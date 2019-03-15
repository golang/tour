// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"net/http"

	"golang.org/x/tools/imports"
)

func init() {
	http.HandleFunc("/fmt", fmtHandler)
}

type fmtResponse struct {
	Body  string
	Error string
}

func fmtHandler(w http.ResponseWriter, r *http.Request) {
	resp := new(fmtResponse)
	var body string
	var err error
	if r.FormValue("imports") == "true" {
		var b []byte
		b, err = imports.Process("prog.go", []byte(r.FormValue("body")), nil)
		body = string(b)
	} else {
		body, err = gofmt(r.FormValue("body"))
	}
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Body = body
	}
	json.NewEncoder(w).Encode(resp)
}

func gofmt(body string) (string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "prog.go", body, parser.ParseComments)
	if err != nil {
		return "", err
	}
	ast.SortImports(fset, f)
	var buf bytes.Buffer
	config := &printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}
	err = config.Fprint(&buf, fset, f)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
