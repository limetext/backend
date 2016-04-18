// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

type (
	// Defines the functionality each package needs to implement
	// so the backend could manage the loading watching and etc
	Package interface {
		Load()
		Name() string
		Path() string
	}

	// We will register each package as a record, Check function for
	// checking if the path suits for the registered package an Action
	// function for creating package from the path
	Record struct {
		Check  func(string) bool
		Action func(string) Package
	}
)

// Keep track of all registered records
var recs []*Record

func Register(r *Record) {
	recs = append(recs, r)
}

func Unregister(r *Record) {
	for i, rec := range recs {
		if rec == r {
			recs, recs[len(recs)-1] = append(recs[:i], recs[i+1:]...), nil
			break
		}
	}
}

func record(path string) Package {
	for _, rec := range recs {
		if !rec.Check(path) {
			continue
		}
		return rec.Action(path)
	}
	return nil
}
