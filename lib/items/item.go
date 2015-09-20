// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package items

import (
	"io/ioutil"
	"path"

	"github.com/limetext/lime-backend/lib/log"
)

type (
	Item interface {
		Load()
		Name() string
	}

	Record struct {
		Check  func(string) bool
		Action func(string) Item
	}
)

var recs []Record

func Register(r Record) {
	recs = append(recs, r)
}

func Scan(dir string) {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Warn("Couldn't read path %s: %s", dir, err)
	}

	watchDir(&pkgDir{dir})

	for _, fi := range fis {
		record(path.Join(dir, fi.Name()))
	}
}

func record(fn string) {
	for _, rec := range recs {
		if rec.Check(fn) {
			item := rec.Action(fn)
			go item.Load()
			Watch(item)
			break
		}
	}
}

func Init() {
	initWatcher()
}
