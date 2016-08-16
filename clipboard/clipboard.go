// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package clipboard

import (
	"github.com/atotto/clipboard"

	"github.com/limetext/backend/log"
)

type (
	Clipboard interface {
		// Get returns the text stored on the clipboard as well as whether or
		// not it was created from an auto-expanded cursor.
		Get() (text string, autoExpanded bool)

		// Set stores text on the clipboard as well as whether or not that text
		// was created from an auto-expanded cursor.
		Set(text string, autoExpanded bool)
	}

	SystemClipboard struct {
		// cachedText is a local copy in case the system clipboard
		// isn't working.
		cachedText string

		// autoExpanded tracks whether the contents was created from a single
		// cursor expanded to a line, by a copy command, for example.
		autoExpanded bool
	}
)

func NewSystemClipboard() *SystemClipboard {
	return &SystemClipboard{}
}

func (c *SystemClipboard) Get() (text string, autoExpanded bool) {
	var err error

	if text, err = clipboard.ReadAll(); err != nil {
		log.Warn("Could not get system clipboard: %v", err)
		text = c.cachedText
	}

	if text == c.cachedText {
		autoExpanded = c.autoExpanded
	}

	return
}

func (c *SystemClipboard) Set(text string, autoExpanded bool) {
	if err := clipboard.WriteAll(text); err != nil {
		log.Warn("Could not set system clipboard: %v", err)
	}

	c.cachedText = text
	c.autoExpanded = autoExpanded
}
