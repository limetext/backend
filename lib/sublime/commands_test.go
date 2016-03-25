// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"reflect"
	"testing"

	"github.com/limetext/lime-backend/lib"
)

func TestCommandGlueInit(t *testing.T) {
	var c CommandGlue

	args := backend.Args{
		"testarg": true,
	}

	if err := c.Init(args); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(c.args, args) {
		t.Error("Expected args to be set, but they weren't")
	}
}

func TestCommandGlueBypassUndo(t *testing.T) {
	var c CommandGlue

	if c.BypassUndo() {
		t.Error("Expected false, but got true")
	}
}
