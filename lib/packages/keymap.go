// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"encoding/json"
)

type Keymap struct {
	simple
}

func NewKeymap(filename string, marshal json.Unmarshaler) *Keymap {
	return &Keymap{simple{filename: filename, marshal: marshal}}
}

func NewKeymapL(filename string, marshal json.Unmarshaler) *Keymap {
	k := NewKeymap(filename, marshal)
	k.Load()
	Watch(k)
	return k
}

// TODO(.): add actions for other events like delete
func (k *Keymap) FileChanged(name string) {
	k.Load()
}
