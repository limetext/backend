// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"io/ioutil"
	"path"

	"github.com/limetext/lime-backend/lib/log"
)

// A helper struct to implement File*Callback interfaces and
// watching all scaned directories for new packages
type scanDir struct {
	path string
}

func (p *scanDir) FileCreated(name string) {
	record(name)
}

// watches scaned directory
func watchDir(dir string) {
	sd := &scanDir{dir}
	if err := watcher.Watch(sd.path, sd); err != nil {
		log.Error("Couldn't watch %s: %s", sd.path, err)
	}
}

func Scan(dir string) error {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	watchDir(dir)

	for _, fi := range fis {
		record(path.Join(dir, fi.Name()))
	}
	return nil
}
