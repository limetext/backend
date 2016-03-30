// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"path"
	"testing"

	"github.com/limetext/lime-backend/lib/keys"
	"github.com/limetext/lime-backend/lib/packages"
)

func TestGetEditor(t *testing.T) {
	editor := GetEditor()
	if editor == nil {
		t.Error("Expected an editor, but got nil")
	}
}

func TestLoadKeyBinding(t *testing.T) {
	ed := GetEditor()
	packages.LoadJSON("testdata/Default.sublime-keymap", ed.KeyBindings())

	kb := ed.KeyBindings().Filter(keys.KeyPress{Key: 'i'})
	if expectedLen := 3; kb.Len() != expectedLen {
		t.Errorf("Expected to have %d keys in the filter, but it had %d", expectedLen, kb.Len())
	}
}

func TestLoadKeyBindings(t *testing.T) {
	ed := GetEditor()
	ed.loadKeyBindings()

	if ed.defaultKB.KeyBindings().Len() <= 0 {
		t.Errorf("Expected editor to have some keys bound, but it didn't")
	}
}

func TestLoadSetting(t *testing.T) {
	editor := GetEditor()
	packages.LoadJSON("testdata/Default.sublime-settings", editor.Settings())

	if editor.Settings().Has("tab_size") != true {
		t.Error("Expected editor settings to have tab_size, but it didn't")
	}

	tab_size := editor.Settings().Get("tab_size").(float64)
	if tab_size != 4 {
		t.Errorf("Expected tab_size to equal 4, got: %v", tab_size)
	}
}

func TestLoadSettings(t *testing.T) {
	editor := GetEditor()
	editor.loadSettings()

	if editor.Settings().Has("tab_size") != true {
		t.Error("Expected editor settings to have tab_size, but it didn't")
	}

	plat := editor.Settings().Parent()
	switch editor.Platform() {
	case "windows":
		if plat.Settings().Get("font_face", "") != "Consolas" {
			t.Errorf("Expected windows font_face be Consolas, but is %s", plat.Settings().Get("font_face", ""))
		}
	case "darwin":
		if plat.Settings().Get("font_face", "") != "Menlo" {
			t.Errorf("Expected OSX font_face be Menlo, but is %s", plat.Settings().Get("font_face", ""))
		}
	default:
		if plat.Settings().Get("font_face", "") != "Monospace" {
			t.Errorf("Expected Linux font_face be Monospace, but is %s", plat.Settings().Get("font_face", ""))
		}
	}
}

func TestInit(t *testing.T) {
	editor := GetEditor()
	editor.Init()

	if editor.defaultKB.KeyBindings().Len() <= 0 {
		t.Errorf("Expected editor to have some keys bound, but it didn't")
	}

	if editor.Settings().Parent().Settings().Parent().Settings().Has("tab_size") != true {
		t.Error("Expected editor settings to have tab_size, but it didn't")
	}
}

func TestNewWindow(t *testing.T) {
	editor := GetEditor()
	l := len(editor.Windows())

	w := editor.NewWindow()
	defer w.Close()

	if len(editor.Windows()) != l+1 {
		t.Errorf("Expected 1 window, but got %d", len(editor.Windows()))
	}
}

func TestRemoveWindow(t *testing.T) {
	editor := GetEditor()
	l := len(editor.Windows())

	w0 := editor.NewWindow()
	defer w0.Close()

	editor.remove(w0)

	if len(editor.Windows()) != l {
		t.Errorf("Expected the window to be removed, but %d still remain", len(editor.Windows()))
	}

	w1 := editor.NewWindow()
	defer w1.Close()

	w2 := editor.NewWindow()
	defer w2.Close()

	editor.remove(w1)

	if len(editor.Windows()) != l+1 {
		t.Errorf("Expected the window to be removed, but %d still remain", len(editor.Windows()))
	}
}

func TestSetActiveWindow(t *testing.T) {
	editor := GetEditor()

	w1 := editor.NewWindow()
	defer w1.Close()

	w2 := editor.NewWindow()
	defer w2.Close()

	if editor.ActiveWindow() != w2 {
		t.Error("Expected the newest window to be active, but it wasn't")
	}

	editor.SetActiveWindow(w1)

	if editor.ActiveWindow() != w1 {
		t.Error("Expected the first window to be active, but it wasn't")
	}
}

func TestSetFrontend(t *testing.T) {
	f := DummyFrontend{}

	editor := GetEditor()
	editor.SetFrontend(&f)

	if editor.Frontend() != &f {
		t.Errorf("Expected a DummyFrontend to be set, but got %T", editor.Frontend())
	}
}

func TestClipboard(t *testing.T) {
	editor := GetEditor()

	// Put back whatever was already there.
	clip := editor.GetClipboard()
	defer editor.SetClipboard(clip)

	want := "test0"

	editor.SetClipboard(want)

	if got := editor.GetClipboard(); got != want {
		t.Errorf("Expected %q to be on the clipboard, but got %q", want, got)
	}

	want = "test1"

	editor.SetClipboard(want)

	if got := editor.GetClipboard(); got != want {
		t.Errorf("Expected %q to be on the clipboard, but got %q", want, got)
	}
}

func TestHandleInput(t *testing.T) {
	// FIXME: This test causes a panic.
	t.Skip("Avoiding pointer issues causing a panic.")

	editor := GetEditor()
	kp := keys.KeyPress{Key: 'i'}

	editor.HandleInput(kp)

	if ki := <-editor.keyInput; ki != kp {
		t.Errorf("Expected %s to be on the input buffer, but got %s", kp, ki)
	}
}

func TestAddColorScheme(t *testing.T) {

}

func TestGetColorScheme(t *testing.T) {

}

func TestColorSchemes(t *testing.T) {

}

func TestAddSyntax(t *testing.T) {

}

func TestGetSyntax(t *testing.T) {

}

func TestSyntaxes(t *testing.T) {

}

func init() {
	ed := GetEditor()
	ed.AddPackagesPath("shipped", path.Join("testdata", "shipped"))
	ed.AddPackagesPath("default", path.Join("testdata", "default"))
	ed.AddPackagesPath("user", path.Join("testdata", "user"))
}
