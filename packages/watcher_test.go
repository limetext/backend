// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatchDir(t *testing.T) {
	// The backend libraries expect absolute paths
	path, err := filepath.Abs("testdata/file")
	if err != nil {
		t.Fatalf("Couldn't get absolute path for %s: %s", path, err)
	}

	pkg := &dummyPackage{path: path}
	rec := &Record{func(s string) bool { return s == path },
		func(s string) Package { return pkg }}

	Register(rec)
	defer Unregister(rec)
	watchDir(filepath.Dir(path))

	if _, err := os.Create(path); err != nil {
		t.Fatalf("Error creating '%s' file: %s", path, err)
	}
	defer os.Remove(path)
	time.Sleep(100 * time.Millisecond)
	if !pkg.IsLoaded() {
		t.Error("Expected package to be loaded")
	}
}
