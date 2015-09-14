// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package items

import (
	"encoding/json"

	"github.com/limetext/lime-backend/lib/log"
)

type Keymap struct {
	simple
}

func NewKeymap(filename string, marshal json.Unmarshaler) *Keymap {
	k := &Keymap{simple{filename: filename, marshal: marshal}}
	watchItem(k)
	return k
}

func NewKeymapL(filename string, marshal json.Unmarshaler) *Keymap {
	k := NewKeymap(filename, marshal)

	if err := k.Load(); err != nil {
		log.Warn("Failed to load keymap %s: %s", k.Name(), err)
	} else {
		log.Info("Loaded keymap %s", k.Name())
	}

	return k
}

// TODO(.): add actions for other events like delete
func (k *Keymap) FileChanged(name string) {
	k.Load()
}
