// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/limetext/text"
)

func TestLoadUnLoad(t *testing.T) {
	set := text.NewSettings()
	j := NewJSON("testdata/Preferences.sublime-settings", &set)

	j.Load()
	if j.err != nil {
		t.Fatalf("Error on loading json %s: %s", j.Name(), j.err)
	}
	if got, exp := set.Get("font_face").(string), "Monospace"; got != exp {
		t.Errorf("Expected font_face %s, but got %s", exp, got)
	}

	j.UnLoad()
	if j.err != nil {
		t.Fatalf("Error on unloading json %s: %s", j.Name(), j.err)
	}
	if set.Has("font_face") {
		t.Error("Expected setting to be empty but has font_face")
	}
}

func TestWatch(t *testing.T) {
	testFile := "testdata/Preferences.sublime-settings"
	testData := []byte(`{"font_face": "test"}`)

	set := text.NewSettings()
	if err := LoadJSON(testFile, &set); err != nil {
		t.Fatalf("Error LoadJSON: %s", err)
	}

	data, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Error reading %s: %s", testFile, err)
	}
	defer ioutil.WriteFile(testFile, data, 0644)

	if err := ioutil.WriteFile(testFile, testData, 0644); err != nil {
		t.Fatalf("Error writing to file %s: %s", testFile, err)
	}
	time.Sleep(100 * time.Millisecond)
	if got, exp := set.Get("font_face").(string), "test"; got != exp {
		t.Errorf("Expected font_face %s, but got %s", exp, got)
	}
}
