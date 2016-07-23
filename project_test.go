// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/limetext/util"
)

func TestSaveAs(t *testing.T) {
	w := GetEditor().NewWindow()
	defer w.Close()
	p := newProject(w)
	p.AddFolder(".")
	p.SaveAs("testdata/saved_project")
	defer os.Remove("testdata/saved_project")

	d1, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Error on marshaling project to json: %s", err)
	}
	d2, err := ioutil.ReadFile("testdata/saved_project")
	if err != nil {
		t.Fatalf("Error on reading 'testdata/saved_project': %s", err)
	}
	if diff := util.Diff(string(d1), string(d2)); diff != "" {
		t.Errorf("Saved project doesn't match expected result\n%s", diff)
	}
}

func TestAddFolder(t *testing.T) {
	w := GetEditor().NewWindow()
	defer w.Close()
	p := newProject(w)
	p.AddFolder("/test/path")

	if got := len(p.Folders()); got != 1 {
		t.Fatalf("Expected project folders len 1, but got %d", got)
	}
	if got, exp := p.Folders()[0], "/test/path"; got != exp {
		t.Errorf("Expected %s in project folders, but got %s", exp, got)
	}
}

func TestRemoveFolder(t *testing.T) {
	w := GetEditor().NewWindow()
	defer w.Close()
	p := newProject(w)
	p.AddFolder("/test/path")
	p.RemoveFolder("/test/path")

	if got := len(p.Folders()); got != 0 {
		t.Errorf("Expected project folders empty, but got %d", got)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/project")
	if err != nil {
		t.Fatalf("Error reading project file 'testdata/project': %s", err)
	}
	w := GetEditor().NewWindow()
	defer w.Close()
	p := newProject(w)
	if err = p.UnmarshalJSON(data); err != nil {
		t.Fatalf("Error on unmarshaling data to project: %s", err)
	}

	if got, exp := p.Settings().Int("tab_size", 4), 8; got != exp {
		t.Errorf("Expected project settings %d, but got %d", exp, got)
	}
	if got := len(p.folders); got != 2 {
		t.Errorf("Expected 2 folders, but got %d", got)
	}

	f1 := p.Folder("src")
	if f1 == nil {
		t.Fatal("Returned src folder is nil")
	}
	if len(f1.ExcludePatterns) == 0 {
		t.Fatal("src 'folder_exluce_patters' is empty")
	}
	if got, exp := f1.ExcludePatterns[0], "backup"; got != exp {
		t.Errorf("Expected %s in src 'folder_exclude_patterns', but got %s", exp, got)
	}
	if !f1.FollowSymlinks {
		t.Error("Expected src 'follow_symlinks' true but its false")
	}

	f2 := p.Folder("docs")
	if f2 == nil {
		t.Fatal("Returned docs folder is nil")
	}
	if got, exp := f2.Name, "Documentation"; got != exp {
		t.Errorf("Expected %s in docs 'name', but got %s", exp, got)
	}
	if len(f2.FileExcludePatterns) == 0 {
		t.Fatal("src 'file_exclude_patterns' is empty")
	}
	if got, exp := f2.FileExcludePatterns[0], "*.css"; got != exp {
		t.Errorf("Expected %s in docs 'file_exclude_patterns', but got %s", exp, got)
	}
}

func TestMarshalJSON(t *testing.T) {
	w := GetEditor().NewWindow()
	defer w.Close()
	p := newProject(w)
	p.Settings().Set("font_size", 12)
	p.AddFolder("./testdata")

	data, err := p.MarshalJSON()
	if err != nil {
		t.Fatalf("Error on marshaling project to json: %s", err)
	}
	exp := `{
	"folders":
	[
		{
			"path": "./testdata"
		}
	],
	"settings":
	{
		"font_size": 12
	}
}
`
	if diff := util.Diff(exp, string(data)); diff != "" {
		t.Logf("Expected:\n%s\nGot:\n%s", exp, string(data))
		t.Errorf("Marshaled project to json doesn't match expected result\n%s", diff)
	}
}

func TestClose(t *testing.T) {
	w := GetEditor().NewWindow()
	defer w.Close()
	p := newProject(w)

	p.Settings().Set("font_size", 14)
	p.AddFolder("./testdata")

	p.Close()
	if got := len(p.Folders()); got != 0 {
		t.Errorf("Expected project folders after close be empty, but got %d", got)
	}
	if got, exp := p.Settings().Int("font_size", 12), 12; got != exp {
		t.Errorf("Expected project font_size settings after close %d, but got %d", exp, got)
	}
}
