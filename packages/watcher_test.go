// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"os"
	"testing"
	"time"
	"path/filepath"
)

func TestWatchDir(t *testing.T) {
	// The backend libraries expect absolute paths
	TestPath, TestPathErr := filepath.Abs("testdata/file")

	if TestPathErr != nil {
		t.Fatalf("Couldn't get absolute path for testdata/file: %s", TestPathErr)
	}

	pkg := &dummyPackage{path: TestPath}
	rec := &Record{func(s string) bool { return s == TestPath },
		func(s string) Package { return pkg }}

	Register(rec)
	defer Unregister(rec)
	watchDir(filepath.Dir(TestPath))

	if _, err := os.Create(TestPath); err != nil {
		t.Fatalf("Error creating '%s' file: %s", TestPath, err)
	}
	defer os.Remove(TestPath)
	time.Sleep(100 * time.Millisecond)
	if !pkg.IsLoaded() {
		t.Error("Expected package loaded")
	}
}
