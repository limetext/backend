// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"sync"
	"testing"
)

type dummyPackage struct {
	path    string
	loaded  bool
	watcher bool
	sync.Mutex
}

func (d *dummyPackage) Load() {
	d.Lock()
	defer d.Unlock()
	d.loaded = true
}
func (d *dummyPackage) UnLoad() {
	d.Lock()
	defer d.Unlock()
	d.loaded = false
}
func (d *dummyPackage) IsLoaded() bool {
	d.Lock()
	defer d.Unlock()
	return d.loaded
}
func (d *dummyPackage) Name() string           { return d.path }
func (d *dummyPackage) Path() string           { return d.path }
func (d *dummyPackage) FileChanged(str string) { d.watcher = true }

func TestRecordCheckAction(t *testing.T) {
	count := 0
	paths := []string{"a", "b", "c", "d"}
	rec := &Record{
		func(s string) bool {
			return (s == "a" || s == "b")
		},
		func(s string) Package {
			count++
			return &dummyPackage{}
		},
	}

	Register(rec)
	defer Unregister(rec)
	for _, path := range paths {
		record(path)
	}
	if count != 2 {
		t.Errorf("Expected count 2 but got: %d", count)
	}
}

func TestRegisterUnregister(t *testing.T) {
	recs = nil

	r1 := &Record{func(s string) bool { return true },
		func(s string) Package { return &dummyPackage{} }}
	r2 := &Record{func(s string) bool { return true },
		func(s string) Package { return &dummyPackage{} }}

	Register(r1)
	if len(recs) != 1 {
		t.Errorf("Expected len of records be 1, but got: %d", len(recs))
	}

	Register(r2)
	if len(recs) != 2 {
		t.Errorf("Expected len of records be 2, but got: %d", len(recs))
	}

	Unregister(r1)
	if len(recs) != 1 {
		t.Errorf("Expected len of records be 1, but got: %d", len(recs))
	}

	Unregister(r2)
	if len(recs) != 0 {
		t.Errorf("Expected len of records be 0, but got: %d", len(recs))
	}
}

func TestLoadUnLoadPackage(t *testing.T) {
	pkg := &dummyPackage{path: "test"}

	load(pkg)
	if !pkg.IsLoaded() {
		t.Error("Expected package be loaded")
	}
	if _, ok := loaded["test"]; !ok {
		t.Error("Expected 'test' in loaded packages")
	}

	unLoad(pkg)
	if pkg.IsLoaded() {
		t.Error("Expected package be unloaded")
	}
	if _, ok := loaded["test"]; ok {
		t.Error("Didn't expect 'test' in loaded packages")
	}
}

func TestUnLoad(t *testing.T) {
	pkg := &dummyPackage{path: "test"}
	load(pkg)

	UnLoad(pkg.Name())
	if pkg.IsLoaded() {
		t.Error("Expected package be unloaded")
	}
}

func TestScan(t *testing.T) {
	path := "testdata/Preferences.sublime-settings"
	pkg := &dummyPackage{path: path}
	rec := &Record{func(s string) bool { return s == path },
		func(s string) Package { return pkg }}

	Register(rec)
	defer Unregister(rec)

	loaded[path] = nil
	Scan("testdata")
	if pkg.IsLoaded() {
		t.Error("Expected not loading package when it exists in loaded packages")
	}

	delete(loaded, path)
	Scan("testdata")
	if !pkg.IsLoaded() {
		t.Error("Expected package be loaded")
	}
}
