// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package items

import (
	"encoding/json"

	"github.com/limetext/lime-backend/lib/log"
)

type Setting struct {
	simple
}

func NewSetting(filename string, marshal json.Unmarshaler) *Setting {
	s := &Setting{simple{filename, marshal}}
	watchItem(s)
	return s
}

func NewSettingL(filename string, marshal json.Unmarshaler) *Setting {
	s := NewSetting(filename, marshal)

	if err := s.Load(); err != nil {
		log.Warn("Failed to load setting %s: %s", s.Name(), err)
	} else {
		log.Info("Loaded setting %s", s.Name())
	}

	return s
}

// TODO(.): add actions for other events like delete
func (s *Setting) FileChanged(name string) {
	s.Load()
}
