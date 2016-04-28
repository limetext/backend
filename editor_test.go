// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"path"
	"sync"
	"testing"

	"github.com/limetext/backend/keys"
	"github.com/limetext/backend/log"
	"github.com/limetext/backend/parser"
	"github.com/limetext/backend/render"
	"github.com/limetext/text"
	qp "github.com/quarnster/parser"
)

type DummyFrontend struct {
	m sync.Mutex
	// Default return value for OkCancelDialog
	defaultAction bool
}

func (h *DummyFrontend) SetDefaultAction(action bool) {
	h.m.Lock()
	defer h.m.Unlock()
	h.defaultAction = action
}
func (h *DummyFrontend) StatusMessage(msg string) { log.Info(msg) }
func (h *DummyFrontend) ErrorMessage(msg string)  { log.Error(msg) }
func (h *DummyFrontend) MessageDialog(msg string) { log.Info(msg) }
func (h *DummyFrontend) OkCancelDialog(msg string, button string) bool {
	log.Info(msg)
	h.m.Lock()
	defer h.m.Unlock()
	return h.defaultAction
}
func (h *DummyFrontend) Show(v *View, r text.Region)       {}
func (h *DummyFrontend) VisibleRegion(v *View) text.Region { return text.Region{} }

func TestGetEditor(t *testing.T) {
	ed := GetEditor()
	if ed == nil {
		t.Error("Expected an editor, but got nil")
	}
}

func TestLoadKeyBindings(t *testing.T) {
	ed := GetEditor()

	if ed.defaultKB.KeyBindings().Len() <= 0 {
		t.Errorf("Expected editor to have some keys bound, but it didn't")
	}
}

func TestLoadSettings(t *testing.T) {
	ed := GetEditor()
	switch ed.Platform() {
	case "windows":
		if res := ed.Settings().Get("font_face", ""); res != "Consolas" {
			t.Errorf("Expected windows font_face be Consolas, but is %s", res)
		}
	case "darwin":
		if res := ed.Settings().Get("font_face", ""); res != "Menlo" {
			t.Errorf("Expected OSX font_face be Menlo, but is %s", res)
		}
	default:
		if res := ed.Settings().Get("font_face", ""); res != "Monospace" {
			t.Errorf("Expected Linux font_face be Monospace, but is %s", res)
		}
	}
}

func TestNewWindow(t *testing.T) {
	ed := GetEditor()
	l := len(ed.Windows())

	w := ed.NewWindow()
	defer w.Close()

	if len(ed.Windows()) != l+1 {
		t.Errorf("Expected 1 window, but got %d", len(ed.Windows()))
	}
}

func TestRemoveWindow(t *testing.T) {
	ed := GetEditor()
	l := len(ed.Windows())

	w0 := ed.NewWindow()
	ed.remove(w0)
	if len(ed.Windows()) != l {
		t.Errorf("Expected the window to be removed, but %d still remain", len(ed.Windows()))
	}

	w1 := ed.NewWindow()
	w2 := ed.NewWindow()
	defer w2.Close()
	ed.remove(w1)
	if len(ed.Windows()) != l+1 {
		t.Errorf("Expected the window to be removed, but %d still remain", len(ed.Windows()))
	}
}

func TestSetActiveWindow(t *testing.T) {
	ed := GetEditor()

	w1 := ed.NewWindow()
	defer w1.Close()

	w2 := ed.NewWindow()
	defer w2.Close()

	if ed.ActiveWindow() != w2 {
		t.Error("Expected the newest window to be active, but it wasn't")
	}

	ed.SetActiveWindow(w1)

	if ed.ActiveWindow() != w1 {
		t.Error("Expected the first window to be active, but it wasn't")
	}
}

func TestSetFrontend(t *testing.T) {
	f := DummyFrontend{}

	ed := GetEditor()
	ed.SetFrontend(&f)

	if ed.Frontend() != &f {
		t.Errorf("Expected a DummyFrontend to be set, but got %T", ed.Frontend())
	}
}

func TestClipboard(t *testing.T) {
	ed := GetEditor()

	// Put back whatever was already there.
	clip := ed.GetClipboard()
	defer ed.SetClipboard(clip)

	want := "test0"

	ed.SetClipboard(want)

	if got := ed.GetClipboard(); got != want {
		t.Errorf("Expected %q to be on the clipboard, but got %q", want, got)
	}

	want = "test1"

	ed.SetClipboard(want)

	if got := ed.GetClipboard(); got != want {
		t.Errorf("Expected %q to be on the clipboard, but got %q", want, got)
	}
}

func TestHandleInput(t *testing.T) {
	// FIXME: This test causes a panic.
	t.Skip("Avoiding pointer issues causing a panic.")

	ed := GetEditor()
	kp := keys.KeyPress{Key: 'i'}

	ed.HandleInput(kp)

	if ki := <-ed.keyInput; ki != kp {
		t.Errorf("Expected %s to be on the input buffer, but got %s", kp, ki)
	}
}

type dummyColorSc struct {
	name string
}

func (d *dummyColorSc) Name() string {
	return d.name
}

func (d *dummyColorSc) Spice(*render.ViewRegions) render.Flavour { return render.Flavour{} }
func (d *dummyColorSc) Global() render.Global                    { return render.Global{} }

func TestAddColorScheme(t *testing.T) {
	cs := new(dummyColorSc)
	ed := GetEditor()

	ed.AddColorScheme("test/path", cs)
	if ret := ed.colorSchemes["test/path"]; ret != cs {
		t.Errorf("Expected 'test/path' color scheme %v, but got %v", cs, ret)
	}
}

type dummySyntax struct {
	name      string
	filetypes []string
	data      string
}

func (d *dummySyntax) Name() string {
	return d.name
}

func (d *dummySyntax) FileTypes() []string {
	return d.filetypes
}

func (d *dummySyntax) Parser(data string) (parser.Parser, error) {
	d.data = data
	return d, nil
}

func (d *dummySyntax) Parse() (*qp.Node, error) { return nil, nil }

func TestAddSyntax(t *testing.T) {
	syn := new(dummySyntax)
	ed := GetEditor()

	ed.AddSyntax("test/path", syn)
	if ret := ed.syntaxes["test/path"]; ret != syn {
		t.Errorf("Expected 'test/path' syntax %v, but got %v", syn, ret)
	}
}

func TestPackagesPath(t *testing.T) {
	ed := GetEditor()
	if got, exp := ed.PackagesPath(), ed.pkgsPaths[0]; exp != got {
		t.Errorf("Expected PackagesPath %s, but got %s", exp, ed.PackagesPath())
	}

	tmp := make([]string, len(ed.pkgsPaths))
	copy(tmp, ed.pkgsPaths)
	ed.pkgsPaths = nil
	if got := ed.PackagesPath(); got != "" {
		t.Errorf("Expected PackagesPath return empty, but got %s", got)
	}
	ed.pkgsPaths = append(ed.pkgsPaths, tmp...)
}

func TestAddRemovePackagesPath(t *testing.T) {
	ed := GetEditor()
	ed.AddPackagesPath("testdata/test")
	l := len(ed.pkgsPaths)

	ed.RemovePackagesPath("testdata/test")
	if got, exp := len(ed.pkgsPaths), l-1; got != exp {
		t.Errorf("Expected len of packages paths %d, but got %d", exp, got)
	}
}

func init() {
	ed := GetEditor()
	ed.Init()
	ed.AddPackagesPath(path.Join("testdata", "Packages"))
	ed.SetDefaultPath(path.Join("testdata", "Packages", "Default"))
	ed.SetUserPath(path.Join("testdata", "Packages", "User"))
}
