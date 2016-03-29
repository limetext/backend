// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/limetext/lime-backend/lib/loaders"
	"github.com/limetext/lime-backend/lib/util"
	"github.com/limetext/text"
)

func TestLoadTheme(t *testing.T) {
	type Test struct {
		in  string
		out string
	}
	tests := []Test{
		{
			"testdata/Monokai.tmTheme",
			"testdata/Monokai.tmTheme.res",
		},
	}
	for _, test := range tests {
		if d, err := ioutil.ReadFile(test.in); err != nil {
			t.Logf("Couldn't load file %s: %s", test.in, err)
		} else {
			var theme Theme
			if err := loaders.LoadPlist(d, &theme); err != nil {
				t.Error(err)
			} else {
				str := fmt.Sprintf("%s", theme)
				if d, err := ioutil.ReadFile(test.out); err != nil {
					if err := ioutil.WriteFile(test.out, []byte(str), 0644); err != nil {
						t.Error(err)
					}
				} else if diff := util.Diff(string(d), str); diff != "" {
					t.Error(diff)
				}

			}
		}
	}
}

func TestLoadThemeFromPlist(t *testing.T) {
	f := "testdata/Monokai.tmTheme"
	th, err := LoadTheme(f)
	if err != nil {
		t.Errorf("Tried to load %s, but got an error: %v", f, err)
	}

	n := "Monokai"
	if th.Name != n {
		t.Errorf("Tried to load %s, but got %s", f, th)
	}
}

func TestLoadThemeFromNonPlist(t *testing.T) {
	f := "testdata/Monokai.tmTheme.res"
	_, err := LoadTheme(f)
	if err == nil {
		t.Errorf("Tried to load %s, expecting an error, but didn't get one", f)
	}
}

func TestLoadThemeFromMissingFile(t *testing.T) {
	f := "testdata/MissingFile"
	_, err := LoadTheme(f)
	if err == nil {
		t.Errorf("Tried to load %s, expecting an error, but didn't get one", f)
	}
}

func TestViewTransform(t *testing.T) {
	w := GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	sc, err := LoadTheme("testdata/GlitterBomb.tmTheme")
	if err != nil {
		t.Fatal(err)
	}

	d, err := ioutil.ReadFile("view.go")
	if err != nil {
		t.Fatal(err)
	}
	e := v.BeginEdit()
	v.Insert(e, 0, string(d))
	v.EndEdit(e)

	if v.Transform(sc, text.Region{A: 0, B: 100}) != nil {
		t.Error("Expected view.Transform return nil when the syntax isn't set yet")
	}

	v.Settings().Set("syntax", "testdata/Go.tmLanguage")

	time.Sleep(time.Second)
	a := v.Transform(sc, text.Region{A: 0, B: 100}).Transcribe()
	v.Transform(sc, text.Region{A: 100, B: 200}).Transcribe()
	c := v.Transform(sc, text.Region{A: 0, B: 100}).Transcribe()
	if !reflect.DeepEqual(a, c) {
		t.Errorf("not equal:\n%v\n%v", a, c)
	}
}

func BenchmarkViewTransformTranscribe(b *testing.B) {
	b.StopTimer()
	w := GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()

	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	sc, err := LoadTheme("testdata/GlitterBomb.tmTheme")
	if err != nil {
		b.Fatal(err)
	}

	v.Settings().Set("syntax", "testdata/Go.tmLanguage")

	d, err := ioutil.ReadFile("view.go")
	if err != nil {
		b.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	v.Settings().AddOnChange("benchmark", func(key string) {
		if key == "lime.syntax.updated" {
			wg.Done()
		}
	})
	e := v.BeginEdit()
	v.Insert(e, 0, string(d))
	v.EndEdit(e)
	wg.Wait()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v.Transform(sc, text.Region{A: 0, B: v.Size()}).Transcribe()
	}
	fmt.Println(util.Prof.String())
}
