// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import "testing"

type dummyPackage struct{}

func (d *dummyPackage) Load()        {}
func (d *dummyPackage) Name() string { return "" }

func TestRecordCheckAction(t *testing.T) {
	Init()

	count := 0
	paths := []string{"a", "b", "c", "d"}
	rec := Record{
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
