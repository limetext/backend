package sublime

import (
	"path"
	"path/filepath"

	"github.com/limetext/gopy/lib"
	"github.com/limetext/lime-backend/lib"
	"github.com/limetext/lime-backend/lib/items"
	"github.com/limetext/lime-backend/lib/keys"
	"github.com/limetext/lime-backend/lib/log"
	"github.com/limetext/text"
)

type plugin struct {
	text.HasSettings
	keys.HasKeyBindings
	platformSet *text.HasSettings
	defaultSet  *text.HasSettings
	defaultKB   *keys.HasKeyBindings
	filename    string
}

func newPlugin(fn string) items.Item {
	p := &plugin{
		filename:    fn,
		platformSet: new(text.HasSettings),
		defaultSet:  new(text.HasSettings),
		defaultKB:   new(keys.HasKeyBindings),
	}
	p.loadKeyBindings()
	p.loadSettings()
	return p
}

func (p *plugin) Load() error {
	fn := p.Name()
	s, err := py.NewUnicode(filepath.Dir(fn) + "." + fn[:len(fn)-3])
	if err != nil {
		return err
	}
	if r, err := module.Base().CallMethodObjArgs("reload_plugin", s); err != nil {
		return err
	} else if r != nil {
		r.Decref()
	}

	return nil
}

func (p *plugin) Name() string {
	return p.filename
}

func (p *plugin) FileChanged(name string) {
	p.Load()
}

func (p *plugin) loadKeyBindings() {
	ed := backend.GetEditor()
	tmp := ed.KeyBindings().Parent()

	ed.KeyBindings().SetParent(p)
	p.KeyBindings().SetParent(p.defaultKB)
	p.defaultKB.KeyBindings().SetParent(tmp)

	pt := path.Join(p.Name(), "Default.sublime-keymap")
	items.NewKeymapL(pt, p.defaultKB.KeyBindings())

	pt = path.Join(p.Name(), "Default ("+ed.Plat()+").sublime-keymap")
	items.NewKeymapL(pt, p.KeyBindings())
}

func (p *plugin) loadSettings() {
	ed := backend.GetEditor()
	tmp := ed.Settings().Parent()

	ed.Settings().SetParent(p)
	p.Settings().SetParent(p.platformSet)
	p.platformSet.Settings().SetParent(p.defaultSet)
	p.defaultSet.Settings().SetParent(tmp)

	pt := path.Join(p.Name(), "Preferences.sublime-settings")
	items.NewSettingL(pt, p.defaultSet.Settings())

	pt = path.Join(p.Name(), "Preferences ("+ed.Plat()+").sublime-settings")
	items.NewSettingL(pt, p.platformSet.Settings())

	pt = path.Join(ed.PackagesPath("user"), "Preferences.sublime-settings")
	items.NewSettingL(pt, p.Settings())
}

func isPlugin(filename string) bool {
	return filepath.Ext(filename) == ".py"
}

func init() {
	backend.OnInit.Add(onInit)
	items.Register(items.Record{isPlugin, newPlugin})
}

var module *py.Module

func onInit() {
	l := py.NewLock()
	defer l.Unlock()

	var err error
	if module, err = py.Import("sublime_plugin"); err != nil {
		panic(err)
	}
	if sys, err := py.Import("sys"); err != nil {
		log.Debug(err)
	} else {
		defer sys.Decref()
	}
}
