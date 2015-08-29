// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package items

import (
	"github.com/limetext/lime-backend/lib/log"
	"github.com/limetext/lime-backend/lib/watch"
)

var watcher *watch.Watcher

func watchItem(item Item) {
	if err := watcher.Watch(item.Name(), item); err != nil {
		log.Warn("Couldn't watch %s: %s", item.Name(), err)
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
