// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"github.com/limetext/lime-backend/lib/log"
	"github.com/limetext/lime-backend/lib/watch"
)

type pkgDir struct {
	path string
}

func (p *pkgDir) FileCreated(name string) {
	record(name)
}

var watcher *watch.Watcher

func Watch(pkg Package) {
	if err := watcher.Watch(pkg.Name(), pkg); err != nil {
		log.Warn("Couldn't watch %s: %s", pkg.Name(), err)
	}
}

func watchDir(p *pkgDir) {
	if err := watcher.Watch(p.path, p); err != nil {
		log.Warn("Couldn't watch %s: %s", p.path, err)
	}
}

func initWatcher() {
	if watcher != nil {
		return
	}

	var err error
	if watcher, err = watch.NewWatcher(); err != nil {
		log.Warn("Couldn't create watcher: %s", err)
	}

	go watcher.Observe()
}
