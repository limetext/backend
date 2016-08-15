// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package clipboard

import (
	"fmt"

	"github.com/limetext/backend/log"
)

type Clipboard struct {
	getter     func() (string, error)
	setter     func(string) error
	cachedText string

	// autoExpanded tracks whether the contents was created from a single
	// cursor expanded to a line, by a copy command, for example.
	autoExpanded bool
}

func New() *Clipboard {
	return &Clipboard{
		getter: func() (string, error) {
			return "", fmt.Errorf("Getter has not been set")
		},
		setter: func(s string) error {
			return fmt.Errorf("Setter has not been set")
		},
	}
}

func (c *Clipboard) SetGetter(getter func() (string, error)) {
	c.getter = getter
}

func (c *Clipboard) SetSetter(setter func(string) error) {
	c.setter = setter
}

func (c *Clipboard) Set(text string, autoExpanded bool) {
	if err := c.setter(text); err != nil {
		log.Warn("Could not set system clipboard: %v", err)
	}

	// Keep a local copy in case the system clipboard isn't working.
	c.cachedText = text
	c.autoExpanded = autoExpanded
}

func (c *Clipboard) Get() (text string, autoExpanded bool) {
	var err error

	if text, err = c.getter(); err != nil {
		log.Warn("Could not get system clipboard: %v", err)
		text = c.cachedText
	}

	if text == c.cachedText {
		autoExpanded = c.autoExpanded
	}

	return
}
