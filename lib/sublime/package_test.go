// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"path"
	"testing"

	_ "github.com/limetext/lime-backend/lib/commands"
	"github.com/limetext/lime-backend/lib/packages"
	_ "github.com/limetext/lime-backend/lib/sublime/api"
)

var (
	pluginPath = path.Join("testdata", "package", "plugin.py")
	pkgPath    = path.Join("testdata", "package")
)

func TestLoadPlugin(t *testing.T) {
	pkg := newPKG(pkgPath).(*pkg)
	pkg.loadPlugin(pluginPath)
	if _, exist := pkg.plugins[pluginPath]; !exist {
		t.Errorf("Expected to %s exist in %s package plugins", pluginPath, pkg.Name())
	}
}

func TestLoadPlugins(t *testing.T) {
	pkg := newPKG(pkgPath).(*pkg)
	pkg.loadPlugins()
	if _, exist := pkg.plugins[pluginPath]; !exist {
		t.Errorf("Expected to %s exist in %s package plugins", pluginPath, pkg.Name())
	}
}

func init() {
	packages.Unregister(packageRecord)
}
