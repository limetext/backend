// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

type (
	// Defines the functionality each package needs to implement
	// so the backend could manage the loading watching and etc
	Package interface {
		Load()

		// Returns the path of the package
		Name() string
	}

	Record struct {
		Check  func(string) bool
		Action func(string) Package
	}
)

var recs []Record

func Register(r Record) {
	recs = append(recs, r)
}

func record(fn string) {
	for _, rec := range recs {
		if rec.Check(fn) {
			pkg := rec.Action(fn)
			go pkg.Load()
			Watch(pkg)
			break
		}
	}
}

func Init() {
	initWatcher()
}
