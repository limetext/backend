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
	filename string
}

func newPlugin(fn string) items.Item {
	return &plugin{filename: fn}
}

func (p *plugin) Load() error {
	dir, file := filepath.Split(p.Name())
	s, err := py.NewUnicode(filepath.Base(dir) + "." + file[:len(file)-3])
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
