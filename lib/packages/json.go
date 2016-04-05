// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"encoding/json"
	"io/ioutil"

	"github.com/limetext/lime-backend/lib/loaders"
	"github.com/limetext/lime-backend/lib/log"
)

// Helper type for loading json files(e.g keymaps settings)
type JSON struct {
	path    string
	marshal json.Unmarshaler
}

func NewJSON(path string, marshal json.Unmarshaler) *JSON {
	return &JSON{path: path, marshal: marshal}
}

// Won't return the json type itself just watch & load
func LoadJSON(path string, marshal json.Unmarshaler) error {
	j := NewJSON(path, marshal)
	if err := watcher.Watch(j.path, j); err != nil {
		log.Warn("Couldn't watch %s: %s", j.path, err)
	}
	return j.Load()
}

func (j *JSON) Load() error {
	log.Debug("Loading %s", j.path)
	data, err := ioutil.ReadFile(j.path)
	if err != nil {
		return err
	}
	return loaders.LoadJSON(data, j.marshal)
}

// TODO(.): add actions for other events like delete and create
func (j *JSON) FileChanged(name string) {
	j.Load()
}
