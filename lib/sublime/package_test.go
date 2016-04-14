// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"path/filepath"
	"testing"

	_ "github.com/limetext/commands"
	"github.com/limetext/lime-backend/lib"
	"github.com/limetext/lime-backend/lib/packages"
	_ "github.com/limetext/lime-backend/lib/sublime/api"
)

var (
	pkgPath    = filepath.Join("testdata", "package")
	pluginPath = filepath.Join("testdata", "package", "plugin.py")
	synPath    = filepath.Join(pkgPath, "Go.tmLanguage")
	csPath     = filepath.Join(pkgPath, "Monokai.tmTheme")
)

func TestLoadPlugin(t *testing.T) {
	pkg := newPKG(pkgPath).(*pkg)
	pkg.loadPlugin(pluginPath)
	checkPlugin(pkg, t)
}

func TestLoadPlugins(t *testing.T) {
	pkg := newPKG(pkgPath).(*pkg)
	pkg.loadPlugins()
	checkPlugin(pkg, t)
}

func TestLoadColorScheme(t *testing.T) {
	pkg := newPKG(pkgPath).(*pkg)
	pkg.loadColorScheme(csPath)
	checkColorScheme(pkg, t)
}

func TestLoadSyntax(t *testing.T) {
	pkg := newPKG(pkgPath).(*pkg)
	pkg.loadSyntax(synPath)
	checkSyntax(pkg, t)
}

func checkPlugin(p *pkg, t *testing.T) {
	if _, exist := p.plugins[pluginPath]; !exist {
		t.Errorf("Expected to %s exist in %s package plugins", pluginPath, p.Name())
	}
}

func checkColorScheme(p *pkg, t *testing.T) {
	if _, ok := p.colorSchemes[csPath]; !ok {
		t.Errorf("Expected %s in %s package color schemes", csPath, p.Name())
	}
	if cs := backend.GetEditor().GetColorScheme(csPath); cs == nil {
		t.Errorf("Expected %s from %s package in editor color schemes", csPath, p.Name())
	}
}

func checkSyntax(p *pkg, t *testing.T) {
	if _, ok := p.syntaxes[synPath]; !ok {
		t.Errorf("Expected %s in %s package syntaxes", synPath, p.Name())
	}
	if syn := backend.GetEditor().GetSyntax(synPath); syn == nil {
		t.Errorf("Expected %s from %s package in editor syntaxes", synPath, p.Name())
	}
}

func TestScan(t *testing.T) {
	pkg := newPKG(pkgPath).(*pkg)
	filepath.Walk(pkg.Path(), pkg.scan)
	checkColorScheme(pkg, t)
	checkSyntax(pkg, t)
}

func init() {
	packages.Unregister(packageRecord)
}
