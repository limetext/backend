package sublime

import (
	"os"
	"path"
	"path/filepath"

	"github.com/limetext/lime-backend/lib"
	"github.com/limetext/lime-backend/lib/keys"
	"github.com/limetext/lime-backend/lib/log"
	"github.com/limetext/lime-backend/lib/packages"
	"github.com/limetext/text"
)

type pkg struct {
	dir string
	text.HasSettings
	keys.HasKeyBindings
	platformSet *text.HasSettings
	defaultSet  *text.HasSettings
	defaultKB   *keys.HasKeyBindings
	plugins     []*plugin
	// TODO: themes, snippets, etc more info on iss#71
}

func newPKG(dir string) packages.Package {
	return &pkg{
		dir:         dir,
		platformSet: new(text.HasSettings),
		defaultSet:  new(text.HasSettings),
		defaultKB:   new(keys.HasKeyBindings),
		plugins:     make([]*plugin, 0),
	}
}

func (p *pkg) Load() {
	log.Debug("Loading package %s", p.Name())
	p.loadKeyBindings()
	p.loadSettings()
	for _, plugin := range p.plugins {
		packages.Watch(plugin)
		plugin.Load()
	}
}

func (p *pkg) Name() string {
	return p.dir
}

func (p *pkg) FileChanged(name string) {}

func (p *pkg) loadKeyBindings() {
	ed := backend.GetEditor()
	tmp := ed.KeyBindings().Parent()
	dir := filepath.Dir(p.Name())

	ed.KeyBindings().SetParent(p)
	p.KeyBindings().SetParent(p.defaultKB)
	p.defaultKB.KeyBindings().SetParent(tmp)

	pt := path.Join(dir, "Default.sublime-keymap")
	packages.NewKeymapL(pt, p.defaultKB.KeyBindings())

	pt = path.Join(dir, "Default ("+ed.Plat()+").sublime-keymap")
	packages.NewKeymapL(pt, p.KeyBindings())
}

func (p *pkg) loadSettings() {
	ed := backend.GetEditor()
	tmp := ed.Settings().Parent()
	dir := filepath.Dir(p.Name())

	ed.Settings().SetParent(p)
	p.Settings().SetParent(p.platformSet)
	p.platformSet.Settings().SetParent(p.defaultSet)
	p.defaultSet.Settings().SetParent(tmp)

	pt := path.Join(dir, "Preferences.sublime-settings")
	packages.NewSettingL(pt, p.defaultSet.Settings())

	pt = path.Join(dir, "Preferences ("+ed.Plat()+").sublime-settings")
	packages.NewSettingL(pt, p.platformSet.Settings())

	pt = path.Join(ed.PackagesPath("user"), "Preferences.sublime-settings")
	packages.NewSettingL(pt, p.Settings())
}

func isPKG(dir string) bool {
	fm, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return fm.IsDir()
}

func init() {
	packages.Register(packages.Record{isPKG, newPKG})
}
