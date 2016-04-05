// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"github.com/limetext/lime-backend/lib/log"
	wch "github.com/limetext/lime-backend/lib/watch"
)

// Responsible for watching all loaded packages
var watcher *wch.Watcher

// Helper function for watching a package
func watch(pkg Package) {
	if err := watcher.Watch(pkg.Name(), pkg); err != nil {
		log.Warn("Couldn't watch %s: %s", pkg.Name(), err)
	}
}

func init() {
	var err error
	if watcher, err = wch.NewWatcher(); err != nil {
		log.Warn("Couldn't create watcher: %s", err)
	}

	go watcher.Observe()
}
