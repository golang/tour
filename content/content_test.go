// Copyright 2016 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package content

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Test that all the .go files inside the content file build
// and execute (without checking for output correctness).
// Files that contain the string "// +build no-build" are not built.
// Files that contain the string "// +build no-run" are not executed.
func TestContent(t *testing.T) {
	scratch, err := ioutil.TempDir("", "tour-content-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(scratch)

	err = filepath.Walk(".", func(path string, fi os.FileInfo, err error) error {
		if filepath.Ext(path) != ".go" {
			return nil
		}
		if filepath.Base(path) == "content_test.go" {
			return nil
		}
		if err := testSnippet(t, path, scratch); err != nil {
			t.Errorf("%v: %v", path, err)
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func testSnippet(t *testing.T, path, scratch string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	build := string(bytes.SplitN(b, []byte{'\n'}, 2)[0])
	if !strings.HasPrefix(build, "// +build ") {
		return errors.New("first line is not a +build comment")
	}
	if !strings.Contains(build, "OMIT") {
		return errors.New(`+build comment does not contain "OMIT"`)
	}

	if strings.Contains(build, "no-build") {
		return nil
	}
	bin := filepath.Join(scratch, filepath.Base(path)+".exe")
	out, err := exec.Command("go", "build", "-o", bin, path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("build error: %v\noutput:\n%s", err, out)
	}
	defer os.Remove(bin)

	if strings.Contains(build, "no-run") {
		return nil
	}
	out, err = exec.Command(bin).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\nOutput:\n%s", err, out)
	}
	return nil
}
