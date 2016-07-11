// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package render

import (
	"encoding/json"
	"fmt"
	"image/color"
	"strconv"

	"github.com/limetext/backend/log"
)

// Color scheme global settings
// http://docs.sublimetext.info/en/latest/reference/color_schemes.html#global-settings-ordered-by-type
type Settings struct {
	Foreground                Colour
	Background                Colour
	Caret                     Colour
	LineHighlight             Colour
	BracketContentsForeground Colour
	// TODO:
	// BracketContentsOptions
	BracketsForeground Colour
	BracketsBackground Colour
	// TODO:
	// BracketsOptions
	TagsForeground Colour
	// TODO:
	// TagsOptions
	FindHighlight           Colour
	FindHighlightForeground Colour
	Gutter                  Colour
	GutterForeground        Colour
	Selection               Colour
	SelectionBackground     Colour
	SelectionBorder         Colour
	InactiveSelection       Colour
	Guide                   Colour
	ActiveGuide             Colour
	StackGuide              Colour
	Highlight               Colour
	HighlightForeground     Colour
	Shadow                  Colour
	// TODO:
	// ShadowWidth
}

// Colour represented by a underlying color.RGBA structure
type Colour color.RGBA

func (c Colour) String() string {
	return fmt.Sprintf("0x%02X%02X%02X%02X", c.A, c.R, c.G, c.B)
}

func (c *Colour) UnmarshalJSON(data []byte) error {
	if data[1] != '#' {
		return c.UnmarshalJSONRGB(data)
	}
	i64, err := strconv.ParseInt(string(data[2:len(data)-1]), 16, 64)
	if err != nil {
		log.Warn("Couldn't properly load color from %s: %s", string(data), err)
	}
	c.A = uint8((i64 >> 24) & 0xff)
	c.R = uint8((i64 >> 16) & 0xff)
	c.G = uint8((i64 >> 8) & 0xff)
	c.B = uint8((i64 >> 0) & 0xff)
	return nil
}

func (c *Colour) UnmarshalJSONRGB(data []byte) error {
	var rgb color.RGBA
	if err := json.Unmarshal(data, &rgb); err != nil {
		log.Warn("Error on unmarshaling %s to color.RGBA: %s", string(data), err)
	}
	*c = Colour(rgb)
	return nil
}
