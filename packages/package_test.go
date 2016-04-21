// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import "testing"

type dummyPackage struct{}

func (d *dummyPackage) Load()        {}
func (d *dummyPackage) UnLoad()      {}
func (d *dummyPackage) Name() string { return "" }
func (d *dummyPackage) Path() string { return "" }

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
