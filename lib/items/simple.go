// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package items

import (
	"encoding/json"
	"io/ioutil"

	"github.com/limetext/lime-backend/lib/loaders"
)

type simple struct {
	filename string
	marshal  json.Unmarshaler
}

func (s *simple) Load() error {
	data, err := ioutil.ReadFile(s.filename)
	if err != nil {
		return err
	}

	return loaders.LoadJSON(data, s.marshal)
}

func (s *simple) Name() string {
	return s.filename
}
