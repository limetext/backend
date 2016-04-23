// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"io/ioutil"
	"path"
	"sync"

	"github.com/limetext/backend/log"
)

type (
	// Defines the functionality each package needs to implement
	// so the backend could manage the loading watching and etc
	Package interface {
		Load()
		UnLoad()
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

var (
	// Registered records
	recs []*Record
	recl sync.Mutex
	// Loaded packages
	loaded = make(map[string]Package)
)

func Register(r *Record) {
	recl.Lock()
	defer recl.Unlock()
	recs = append(recs, r)
}

func Unregister(r *Record) {
	recl.Lock()
	defer recl.Unlock()
	for i, rec := range recs {
		if rec == r {
			recs, recs[len(recs)-1] = append(recs[:i], recs[i+1:]...), nil
			break
		}
	}
}

func Scan(dir string) {
	log.Debug("Scanning %s for packages", dir)
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Error("Error while scanning %s: %s", dir, err)
	}
	watchDir(dir)

	var pkgs []Package
	for _, fi := range fis {
		pkgPath := path.Join(dir, fi.Name())
		if _, ok := loaded[pkgPath]; ok {
			continue
		}
		if pkg := record(pkgPath); pkg != nil {
			pkgs = append(pkgs, pkg)
		}
	}
	// TODO: we cant run this in a go routine because currently there is
	// no way to frontends to know when for example the color scheme
	// is ready
	func() {
		for _, pkg := range pkgs {
			load(pkg)
		}
	}()
}

func UnLoad(name string) {
	for _, pkg := range loaded {
		if pkg.Name() == name {
			unLoad(pkg)
			return
		}
	}
}

func record(path string) Package {
	recl.Lock()
	defer recl.Unlock()
	for _, rec := range recs {
		if rec.Check(path) {
			return rec.Action(path)
		}
	}
	return nil
}

func load(pkg Package) {
	pkg.Load()
	watch(pkg)
	loaded[pkg.Path()] = pkg
}

func unLoad(pkg Package) {
	pkg.UnLoad()
	unWatch(pkg)
	delete(loaded, pkg.Path())
}
