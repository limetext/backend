// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/limetext/lime-backend/lib"
	"github.com/limetext/lime-backend/lib/keys"
	"github.com/limetext/lime-backend/lib/log"
	"github.com/limetext/lime-backend/lib/packages"
	_ "github.com/limetext/lime-backend/lib/sublime/api"
	"github.com/limetext/text"
)

// A sublime package
type pkg struct {
	dir  string
	name string
	text.HasSettings
	keys.HasKeyBindings
	platformSettings *text.HasSettings
	defaultSettings  *text.HasSettings
	defaultKB        *keys.HasKeyBindings
	plugins          map[string]*plugin
	syntaxes         map[string]*Syntax
	colorSchemes     map[string]*ColorScheme
	// TODO: themes, snippets, etc more info on iss#71
}

func newPKG(dir string) packages.Package {
	p := &pkg{
		dir:              dir,
		name:             filepath.Base(dir),
		platformSettings: new(text.HasSettings),
		defaultSettings:  new(text.HasSettings),
		defaultKB:        new(keys.HasKeyBindings),
		plugins:          make(map[string]*plugin),
		syntaxes:         make(map[string]*Syntax),
		colorSchemes:     make(map[string]*ColorScheme),
	}

	ed := backend.GetEditor()

	// Initializing settings hierarchy
	// editor <- default <- platform <- user(package)
	p.Settings().SetParent(p.platformSettings)
	p.platformSettings.Settings().SetParent(p.defaultSettings)
	p.defaultSettings.Settings().SetParent(ed)

	// Initializing keybidings hierarchy
	// editor <- default <- platform(package)
	p.KeyBindings().SetParent(p.defaultKB)
	p.defaultKB.KeyBindings().SetParent(ed)

	return p
}

func (p *pkg) Load() {
	log.Debug("Loading package %s", p.Name())
	p.loadKeyBindings()
	p.loadSettings()
	p.loadPlugins()

	filepath.Walk(p.Name(), p.scan)
}

func (p *pkg) Path() string {
	return p.dir
}

func (p *pkg) Name() string {
	return p.name
}

// TODO: how we should watch the package and the files containing?
func (p *pkg) FileCreated(name string) {}

func (p *pkg) loadPlugins() {
	log.Fine("Loading %s plugins", p.Name())
	fis, err := ioutil.ReadDir(p.Path())
	if err != nil {
		log.Warn("Error on reading directory %s, %s", p.Path(), err)
		return
	}
	for _, fi := range fis {
		if isPlugin(fi.Name()) {
			p.loadPlugin(filepath.Join(p.Path(), fi.Name()))
		}
	}
}

func (p *pkg) loadPlugin(path string) {
	pl := newPlugin(path)
	pl.Load()

	p.plugins[path] = pl.(*plugin)
}

func (p *pkg) loadColorScheme(path string) {
	log.Debug("Loading color scheme %s", path)
	tm, err := LoadTheme(path)
	if err != nil {
		log.Warn("Error loading %s color scheme %s: %s", p.Name(), path, err)
		return
	}

	cs := &ColorScheme{*tm}
	// TODO: the path should be modified
	p.colorSchemes[path] = cs
	backend.GetEditor().AddColorScheme(path, cs)
}

func (p *pkg) loadSyntax(path string) {
	log.Debug("Loading syntax %s", path)
	syn, err := newSyntax(path)
	if err != nil {
		log.Warn("Error loading %s syntax: %s", p.Name(), err)
		return
	}

	// TODO: the path should be modified
	p.syntaxes[path] = syn
	backend.GetEditor().AddSyntax(path, syn)
}

func (p *pkg) loadKeyBindings() {
	log.Fine("Loading %s keybindings", p.Name())
	ed := backend.GetEditor()

	pt := filepath.Join(p.Path(), "Default.sublime-keymap")
	log.Finest("Loading %s", pt)
	packages.LoadJSON(pt, p.defaultKB.KeyBindings())

	pt = filepath.Join(p.Path(), "Default ("+ed.Plat()+").sublime-keymap")
	log.Finest("Loading %s", pt)
	packages.LoadJSON(pt, p.KeyBindings())
}

func (p *pkg) loadSettings() {
	log.Fine("Loading %s settings", p.Name())
	ed := backend.GetEditor()

	pt := filepath.Join(p.Path(), "Preferences.sublime-settings")
	log.Finest("Loading %s", pt)
	packages.LoadJSON(pt, p.defaultSettings.Settings())

	pt = filepath.Join(p.Path(), "Preferences ("+ed.Plat()+").sublime-settings")
	log.Finest("Loading %s", pt)
	packages.LoadJSON(pt, p.platformSettings.Settings())

	pt = filepath.Join(ed.PackagesPath("user"), "Preferences.sublime-settings")
	log.Finest("Loading %s", pt)
	packages.LoadJSON(pt, p.Settings())
}

func (p *pkg) scan(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	if isColorScheme(path) {
		p.loadColorScheme(path)
	}
	if isSyntax(path) {
		p.loadSyntax(path)
	}
	return nil
}

func isColorScheme(path string) bool {
	if filepath.Ext(path) == ".tmTheme" {
		return true
	}
	return false
}

func isSyntax(path string) bool {
	if filepath.Ext(path) == ".tmLanguage" {
		return true
	}
	return false
}

// Any directory in sublime is a package
func isPKG(dir string) bool {
	fi, err := os.Stat(dir)
	if err != nil || !fi.IsDir() {
		return false
	}
	return true
}

var packageRecord = &packages.Record{isPKG, newPKG}

func init() {
	packages.Register(packageRecord)
}
