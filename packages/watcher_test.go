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
	testPath, testPathErr := filepath.Abs("testdata/file")

	if testPathErr != nil {
		t.Fatalf("Couldn't get absolute path for testdata/file: %s", testPathErr)
	}

	pkg := &dummyPackage{path: testPath}
	rec := &Record{func(s string) bool { return s == testPath },
		func(s string) Package { return pkg }}

	Register(rec)
	defer Unregister(rec)
	watchDir(filepath.Dir(testPath))

	if _, err := os.Create(testPath); err != nil {
		t.Fatalf("Error creating '%s' file: %s", testPath, err)
	}
	defer os.Remove(testPath)
	time.Sleep(100 * time.Millisecond)
	if !pkg.IsLoaded() {
		t.Error("Expected package loaded")
	}
}
