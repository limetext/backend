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
	filename string
	marshal  json.Unmarshaler
}

func NewJSON(filename string, marshal json.Unmarshaler) *JSON {
	return &JSON{filename: filename, marshal: marshal}
}

// Won't return the json type itself just watch & load
func LoadJSON(filename string, marshal json.Unmarshaler) error {
	j := NewJSON(filename, marshal)
	if err := watcher.Watch(j.Name(), j); err != nil {
		log.Warn("Couldn't watch %s: %s", j.Name(), err)
	}
	return j.Load()
}

func (j *JSON) Load() error {
	log.Debug("Loading %s", j.Name())
	data, err := ioutil.ReadFile(j.filename)
	if err != nil {
		return err
	}
	return loaders.LoadJSON(data, j.marshal)
}

func (j *JSON) Name() string {
	return j.filename
}

// TODO(.): add actions for other events like delete and create
func (j *JSON) FileChanged(name string) {
	j.Load()
}
