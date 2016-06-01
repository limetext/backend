// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"fmt"

	"github.com/limetext/backend/render"
)

// Any color scheme view should implement this interface
// also it should register it self from editor.AddColorSCheme
type ColorScheme interface {
	render.ColourScheme
	Name() string
}

func colorScheme(name string) (ColorScheme, error) {
	if name == "" {
		// TODO bring the default color scheme up
	}
	scheme := ed.GetColorScheme(name)
	if scheme == nil {
		return nil, fmt.Errorf("No color scheme %s in editor", name)
	}
	return scheme, nil
}
