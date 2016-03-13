// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"encoding/json"
)

type Setting struct {
	simple
}

func NewSetting(filename string, marshal json.Unmarshaler) *Setting {
	return &Setting{simple{filename: filename, marshal: marshal}}
}

func NewSettingL(filename string, marshal json.Unmarshaler) *Setting {
	s := NewSetting(filename, marshal)
	s.Load()
	wch(s)
	return s
}

// TODO(.): add actions for other events like delete
func (s *Setting) FileChanged(name string) {
	s.Load()
}
