// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"testing"

	"github.com/limetext/sublime/textmate/theme"
	"github.com/limetext/text"
)

type dummyColorScheme struct {
	*theme.Theme
}

func newDummyColorScheme(tb testing.TB, path string) *dummyColorScheme {
	if tm, err := theme.Load(path); err != nil {
		tb.Fatalf("Error loading theme %s: %s", path, err)
		return nil
	} else {
		return &dummyColorScheme{tm}
	}
}

func (c *dummyColorScheme) Name() string {
	return c.Theme.Name
}

func addSetColorScheme(tb testing.TB, settings *text.Settings, path string) {
	cs := newDummyColorScheme(tb, path)
	GetEditor().AddColorScheme(path, cs)
	settings.Set("colour_scheme", path)
}
