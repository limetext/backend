// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"encoding/json"
	"io/ioutil"

	"github.com/limetext/loaders"
)

// Helper type for loading json files(e.g keymaps settings)
type JSON struct {
	path    string
	err     error
	marshal json.Unmarshaler
}

func NewJSON(path string, marshal json.Unmarshaler) *JSON {
	return &JSON{path: path, marshal: marshal}
}

// Won't return the json type itself just watch & load
func LoadJSON(path string, marshal json.Unmarshaler) error {
	j := NewJSON(path, marshal)
	watch(j)
	j.Load()
	return j.err
}

func (j *JSON) Load() {
	j.err = nil
	data, err := ioutil.ReadFile(j.Path())
	if err != nil {
		j.err = err
		return
	}
	j.err = loaders.LoadJSON(data, j.marshal)
}

func (j *JSON) UnLoad() {
	json.Unmarshal([]byte(``), j.marshal)
}

func (j *JSON) Name() string { return j.path }

func (j *JSON) Path() string { return j.path }

// TODO(.): add actions for other events like delete and create
func (j *JSON) FileChanged(name string) {
	j.Load()
}
