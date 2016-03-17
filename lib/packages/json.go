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

// Helper struct for simple packages containing 1 json file(e.g keymaps settings)
type JSON struct {
	filename string
	marshal  json.Unmarshaler
}

func NewJSON(filename string, marshal json.Unmarshaler) *JSON {
	return &JSON{filename: filename, marshal: marshal}
}

func NewJSONL(filename string, marshal json.Unmarshaler) *JSON {
	j := NewJSON(filename, marshal)
	j.Load()
	wch(j)
	return j
}

// TODO: better errors, maybe we should introduce error type and let
// the load caller decide how to log
func (j *JSON) Load() {
	log.Debug("Loading %s", j.Name())
	data, err := ioutil.ReadFile(j.filename)
	if err != nil {
		log.Warn(err)
		return
	}

	if err = loaders.LoadJSON(data, j.marshal); err != nil {
		log.Warn(err)
	}
}

func (j *JSON) Name() string {
	return j.filename
}
