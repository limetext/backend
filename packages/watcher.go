// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"github.com/limetext/backend/log"
	wch "github.com/limetext/backend/watch"
)

// Responsible for watching all loaded packages
var watcher *wch.Watcher

// Helper function for watching a package
func watch(pkg Package) {
	if err := watcher.Watch(pkg.Path(), pkg); err != nil {
		log.Warn("Couldn't watch %s: %s", pkg.Path(), err)
	}
}

func unWatch(pkg Package) {
	if err := watcher.UnWatch(pkg.Path(), pkg); err != nil {
		log.Warn("Couldnt unwatch %s: %s", pkg.Path(), err)
	}
}

// A helper struct to implement File*Callback interfaces and
// watching all scaned directories for new packages
type scanDir struct {
	path string
}

// TODO: are we checking new folders to?
func (p *scanDir) FileCreated(name string) {
	if pkg := record(name); pkg != nil {
		load(pkg)
	}
}

// watches scaned directory
func watchDir(dir string) {
	log.Finest("Watching scaned dir: %s", dir)
	sd := &scanDir{dir}
	if err := watcher.Watch(sd.path, sd); err != nil {
		log.Error("Couldn't watch %s: %s", sd.path, err)
	}
}

func init() {
	var err error
	if watcher, err = wch.NewWatcher(); err != nil {
		log.Warn("Couldn't create watcher: %s", err)
	}

	go watcher.Observe()
}
