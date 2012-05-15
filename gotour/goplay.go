// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/compile", Compile)
}

type Response struct {
	Output string `json:"output"`
	Errors string `json:"compile_errors"`
}

func Compile(w http.ResponseWriter, req *http.Request) {
	resp := new(Response)
	out, err := compile(req)
	if err != nil {
		if len(out) > 0 {
			resp.Errors = string(out) + "\n" + err.Error()
		} else {
			resp.Errors = err.Error()
		}
	} else {
		resp.Output = string(out)
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println(err)
	}
}
