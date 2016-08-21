// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"os"
	"testing"
	"time"
	"strings"
)

func TestWatchDir(t *testing.T) {
	pkg := &dummyPackage{path: "testdata/file"}
	rec := &Record{func(s string) bool { return strings.Contains(s, "testdata/file") },
		func(s string) Package { return pkg }}

	Register(rec)
	defer Unregister(rec)
	watchDir("testdata")

	if _, err := os.Create("testdata/file"); err != nil {
		t.Fatalf("Error creating 'testdata/file' file: %s", err)
	}
	defer os.Remove("testdata/file")
	time.Sleep(100 * time.Millisecond)
	if !pkg.IsLoaded() {
		t.Error("Expected package loaded")
	}
}
