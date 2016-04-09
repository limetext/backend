// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/limetext/gopy/lib"
	"github.com/limetext/lime-backend/lib"
	_ "github.com/limetext/lime-backend/lib/sublime/api"
)

func TestPlugin(t *testing.T) {
	ed := backend.GetEditor()
	ed.AddPackagesPath("test", path.Join("testdata", "plugins"))
	time.Sleep(time.Millisecond * 100)

	l := py.NewLock()
	defer l.Unlock()
	if _, err := py.Import("plugin_test"); err != nil {
		t.Error(err)
	}
}

func TestReloadPlugin(t *testing.T) {
	data := []byte(`import sublime, sublime_plugin

class TestToxt(sublime_plugin.TextCommand):
    def run(self, edit):
        self.view.insert(edit, 0, "Tada")
		`)
	if err := ioutil.WriteFile("testdata/plugins/reload.py", data, 0644); err != nil {
		t.Fatalf("Couldn't write testdata/plugins/reload.py: %s", err)
	}
	defer os.Remove("testdata/plugins/reload.py")
	time.Sleep(time.Millisecond * 100)

	l := py.NewLock()
	defer l.Unlock()
	if _, err := py.Import("reload_test"); err != nil {
		t.Error(err)
	}
}

func init() {
	pyAddPath(".")
	pyAddPath("testdata")

	ed := backend.GetEditor()
	ed.Init()
	ed.NewWindow()
}
