// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/limetext/lime-backend/lib"
	"github.com/limetext/lime-backend/lib/keys"
	"github.com/limetext/lime-backend/lib/log"
	"github.com/limetext/lime-backend/lib/packages"
	"github.com/limetext/text"
)

// A sublime package
type pkg struct {
	dir string
	text.HasSettings
	keys.HasKeyBindings
	platformSet *text.HasSettings
	defaultSet  *text.HasSettings
	defaultKB   *keys.HasKeyBindings
	plugins     map[string]*plugin
	// TODO: themes, snippets, etc more info on iss#71
}

func newPKG(dir string) packages.Package {
	return &pkg{
		dir:         dir,
		platformSet: new(text.HasSettings),
		defaultSet:  new(text.HasSettings),
		defaultKB:   new(keys.HasKeyBindings),
		plugins:     make(map[string]*plugin),
	}
}

func (p *pkg) Load() {
	log.Debug("Loading package %s", p.Name())
	p.loadKeyBindings()
	p.loadSettings()
	p.loadPlugins()
}

func (p *pkg) Name() string {
	return p.dir
}

// TODO: how we should watch the package and the files containing?
func (p *pkg) FileCreated(name string) {
	p.loadPlugin(name)
}

func (p *pkg) loadPlugins() {
	log.Fine("Loading %s plugins", p.Name())
	fis, err := ioutil.ReadDir(p.Name())
	if err != nil {
		log.Warn("Error on reading directory %s, %s", p.Name(), err)
		return
	}
	for _, fi := range fis {
		if isPlugin(fi.Name()) {
			p.loadPlugin(path.Join(p.Name(), fi.Name()))
		}
	}
}

func (p *pkg) loadPlugin(fn string) {
	if _, exist := p.plugins[fn]; exist {
		return
	}

	pl := newPlugin(fn)
	pl.Load()

	p.plugins[fn] = pl.(*plugin)
}

func (p *pkg) loadKeyBindings() {
	log.Fine("Loading %s keybindings", p.Name())
	ed := backend.GetEditor()
	tmp := ed.KeyBindings().Parent()
	dir := filepath.Dir(p.Name())

	ed.KeyBindings().SetParent(p)
	p.KeyBindings().SetParent(p.defaultKB)
	p.defaultKB.KeyBindings().SetParent(tmp)

	pt := path.Join(dir, "Default.sublime-keymap")
	packages.NewJSONL(pt, p.defaultKB.KeyBindings())

	pt = path.Join(dir, "Default ("+ed.Plat()+").sublime-keymap")
	packages.NewJSONL(pt, p.KeyBindings())
}

func (p *pkg) loadSettings() {
	log.Fine("Loading %s settings", p.Name())
	ed := backend.GetEditor()
	tmp := ed.Settings().Parent()
	dir := filepath.Dir(p.Name())

	ed.Settings().SetParent(p)
	p.Settings().SetParent(p.platformSet)
	p.platformSet.Settings().SetParent(p.defaultSet)
	p.defaultSet.Settings().SetParent(tmp)

	pt := path.Join(dir, "Preferences.sublime-settings")
	packages.NewJSONL(pt, p.defaultSet.Settings())

	pt = path.Join(dir, "Preferences ("+ed.Plat()+").sublime-settings")
	packages.NewJSONL(pt, p.platformSet.Settings())

	pt = path.Join(ed.PackagesPath("user"), "Preferences.sublime-settings")
	packages.NewJSONL(pt, p.Settings())
}

// Any directory in sublime is a package
func isPKG(dir string) bool {
	fi, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

var packageRecord *packages.Record = &packages.Record{isPKG, newPKG}

func init() {
	packages.Register(packageRecord)
}
