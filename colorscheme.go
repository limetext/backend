// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"github.com/limetext/backend/log"
	"github.com/limetext/backend/render"
)

// Any color scheme view should implement this interface
// also it should register it self from editor.AddColorSCheme
type ColorScheme interface {
	render.ColourScheme
	Name() string
}

func colorScheme(name string) render.ColourScheme {
	if name == "" {
		return defaultScheme()
	}

	scheme := ed.GetColorScheme(name)
	if scheme == nil {
		log.Error("No color scheme %s in editor falling back to default color scheme", name)
		return defaultScheme()
	}
	return scheme
}

type scheme struct {
	settings render.Settings
}

func (s *scheme) Spice(*render.ViewRegions) render.Flavour {
	return render.Flavour{
		Background: s.Settings().Background,
		Foreground: s.Settings().Foreground,
	}
}

func (s *scheme) Settings() render.Settings {
	return s.settings
}

// default colorscheme used when there is a problem
var colorscheme *scheme

func defaultScheme() render.ColourScheme {
	if colorscheme == nil {
		colorscheme = &scheme{
			render.Settings{
				Background: render.Colour{255, 255, 255, 1},
			},
		}
	}
	return colorscheme
}
