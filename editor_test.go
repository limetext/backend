// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"path"
	"testing"

	"github.com/limetext/backend/keys"
)

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
		if res := ed.Settings().String("font_face", ""); res != "Consolas" {
			t.Errorf("Expected windows font_face be Consolas, but is %s", res)
		}
	case "darwin":
		if res := ed.Settings().String("font_face", ""); res != "Menlo" {
			t.Errorf("Expected OSX font_face be Menlo, but is %s", res)
		}
	default:
		if res := ed.Settings().String("font_face", ""); res != "Monospace" {
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
	f := dummyFrontend{}

	ed := GetEditor()
	ed.SetFrontend(&f)

	if ed.Frontend() != &f {
		t.Errorf("Expected a dummyFrontend to be set, but got %T", ed.Frontend())
	}
}

func TestClipboard(t *testing.T) {
	ed := GetEditor()
	// Put back whatever was already there.
	clip, ex := ed.GetClipboard()
	defer ed.SetClipboard(clip, ex)

	want := "test0"
	ed.SetClipboard(want, true)

	got, ex := ed.GetClipboard()

	if got != want {
		t.Errorf("Expected %q to be on the clipboard, but got %q", want, got)
	}

	if !ex {
		t.Errorf("Expected the clipboard to be flagged as auto expanded, but it wasn't")
	}

	want = "test1"
	ed.SetClipboard(want, true)

	got, ex = ed.GetClipboard()

	if got != want {
		t.Errorf("Expected %q to be on the clipboard, but got %q", want, got)
	}

	if !ex {
		t.Errorf("Expected the clipboard to be flagged as auto expanded, but it wasn't")
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

func TestAddColorScheme(t *testing.T) {
	csPath := "testdata/Monokai.tmTheme"
	cs := newDummyColorScheme(t, csPath)
	ed := GetEditor()

	ed.AddColorScheme(csPath, cs)
	if ret := ed.colorSchemes[csPath]; ret != cs {
		t.Errorf("Expected '%s' color scheme %v, but got %v", csPath, cs, ret)
	}
}

func TestAddSyntax(t *testing.T) {
	synPath := "testdata/Go.tmLanguage"
	syn := newDummySytax(t, synPath)
	ed := GetEditor()

	ed.AddSyntax(synPath, syn)
	if ret := ed.syntaxes[synPath]; ret != syn {
		t.Errorf("Expected '%s' syntax %v, but got %v", synPath, syn, ret)
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
