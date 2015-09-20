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

type simple struct {
	filename string
	marshal  json.Unmarshaler
}

func (s *simple) Load() {
	log.Debug("Loading %s", s.Name())
	data, err := ioutil.ReadFile(s.filename)
	if err != nil {
		log.Warn(err)
		return
	}

	if err = loaders.LoadJSON(data, s.marshal); err != nil {
		log.Warn(err)
	}
}

func (s *simple) Name() string {
	return s.filename
}
