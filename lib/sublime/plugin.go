package sublime

import (
	"path/filepath"

	"github.com/limetext/gopy/lib"
	"github.com/limetext/lime-backend/lib"
	"github.com/limetext/lime-backend/lib/log"
	"github.com/limetext/lime-backend/lib/packages"
)

type plugin struct {
	filename string
}

func newPlugin(fn string) packages.Package {
	return &plugin{filename: fn}
}

func (p *plugin) Load() {
	log.Debug("Loading plugin %s", p.Name())
	dir, file := filepath.Split(p.Name())
	s, err := py.NewUnicode(filepath.Base(dir) + "." + file[:len(file)-3])
	if err != nil {
		log.Warn(err)
		return
	}

	l := py.NewLock()
	defer l.Unlock()
	if r, err := module.Base().CallMethodObjArgs("reload_plugin", s); err != nil {
		log.Warn(err)
		return
	} else if r != nil {
		r.Decref()
	}
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
	packages.Register(packages.Record{isPlugin, newPlugin})
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
