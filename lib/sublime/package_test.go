// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"path"
	"testing"

	_ "github.com/limetext/lime-backend/lib/commands"
	"github.com/limetext/lime-backend/lib/packages"
)

func TestLoadPlugin(t *testing.T) {
	pn := path.Join("testdata", "packages", "plugin.py")
	pkg := newPKG(path.Join("testdata", "packages")).(*pkg)
	pkg.loadPlugin(pn)

	if _, exist := pkg.plugins[pn]; !exist {
		t.Errorf("Expected to %s exist in %s package plugins", pn, pkg.Name())
	}
}

func TestLoadPlugins(t *testing.T) {
	pn := path.Join("testdata", "packages", "plugin.py")
	pkg := newPKG(path.Join("testdata", "packages")).(*pkg)
	pkg.loadPlugins()

	if _, exist := pkg.plugins[pn]; !exist {
		t.Errorf("Expected to %s exist in %s package plugins", pn, pkg.Name())
	}
}

func init() {
	packages.Unregister(packageRecord)
}
