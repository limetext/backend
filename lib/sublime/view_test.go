// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/limetext/lime-backend/lib"
	"github.com/limetext/lime-backend/lib/util"
	"github.com/limetext/text"
)

func TestViewTransform(t *testing.T) {
	w := backend.GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	sc, err := LoadTheme("testdata/package/Monokai.tmTheme")
	if err != nil {
		t.Fatal(err)
	}

	d, err := ioutil.ReadFile("testdata/view.go")
	if err != nil {
		t.Fatal(err)
	}
	e := v.BeginEdit()
	v.Insert(e, 0, string(d))
	v.EndEdit(e)

	if v.Transform(sc, text.Region{A: 0, B: 100}) != nil {
		t.Error("Expected view.Transform return nil when the syntax isn't set yet")
	}

	v.Settings().Set("syntax", "testdata/package/Go.tmLanguage")

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
	w := backend.GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()

	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	sc, err := LoadTheme("testdata/package/Monokai.tmTheme")
	if err != nil {
		b.Fatal(err)
	}

	v.Settings().Set("syntax", "testdata/package/Go.tmLanguage")

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

// This is not 100% what ST3 does
func TestViewExtractScope(t *testing.T) {
	w := backend.GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	const (
		in      = "testdata/main.go"
		expfile = "testdata/scoperange.res"
		syntax  = "testdata/package/Go.tmLanguage"
	)
	syn, err := newSyntax(syntax)
	if err != nil {
		t.Fatal(err)
	}
	backend.GetEditor().AddSyntax(syntax, syn)
	v.Settings().Set("syntax", syntax)
	d, err := ioutil.ReadFile(in)
	if err != nil {
		t.Fatal(err)
	}
	e := v.BeginEdit()
	v.Insert(e, 0, string(d))
	v.EndEdit(e)
	last := text.Region{A: -1, B: -1}
	str := ""
	nr := text.Region{A: 0, B: 0}
	for v.ExtractScope(1) == nr {
		time.Sleep(time.Millisecond)
	}
	for i := 0; i < v.Size(); i++ {
		if r := v.ExtractScope(i); r != last {
			str += fmt.Sprintf("%d (%d, %d)\n", i, r.A, r.B)
			last = r
		}
	}
	if d, err := ioutil.ReadFile(expfile); err != nil {
		if err := ioutil.WriteFile(expfile, []byte(str), 0644); err != nil {
			t.Error(err)
		}
	} else if diff := util.Diff(string(d), str); diff != "" {
		t.Error(diff)
	}
}

// This is not 100% what ST3 does, but IMO ST3 is wrong
func TestViewScopeName(t *testing.T) {
	w := backend.GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	const (
		in      = "testdata/main.go"
		expfile = "testdata/scopename.res"
		syntax  = "testdata/package/Go.tmLanguage"
	)
	syn, err := newSyntax(syntax)
	if err != nil {
		t.Fatal(err)
	}
	backend.GetEditor().AddSyntax(syntax, syn)
	v.Settings().Set("syntax", syntax)
	d, err := ioutil.ReadFile(in)
	if err != nil {
		t.Fatal(err)
	}
	e := v.BeginEdit()
	v.Insert(e, 0, string(d))
	v.EndEdit(e)
	last := ""
	str := ""
	lasti := 0
	for v.ScopeName(1) == "" {
		time.Sleep(250 * time.Millisecond)
	}
	for i := 0; i < v.Size(); i++ {
		if name := v.ScopeName(i); name != last {
			if last != "" {
				str += fmt.Sprintf("%d-%d: %s\n", lasti, i, last)
				lasti = i
			}
			last = name
		}
	}
	if i := v.Size(); lasti != i {
		str += fmt.Sprintf("%d-%d: %s\n", lasti, i, last)
	}
	if d, err := ioutil.ReadFile(expfile); err != nil {
		if err := ioutil.WriteFile(expfile, []byte(str), 0644); err != nil {
			t.Error(err)
		}
	} else if diff := util.Diff(string(d), str); diff != "" {
		t.Error(diff)
	}
}

func TestViewStress(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	ed := backend.GetEditor()
	w := ed.NewWindow()
	defer w.Close()

	v := w.OpenFile("testdata/view.go", 0)
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	syntax := "testdata/package/Go.tmLanguage"
	syn, err := newSyntax(syntax)
	if err != nil {
		t.Fatal(err)
	}
	backend.GetEditor().AddSyntax(syntax, syn)
	v.Settings().Set("syntax", syntax)
	for i := 0; i < 1000; i++ {
		e := v.BeginEdit()
		for i := 0; i < 100; i++ {
			v.Insert(e, 0, "h")
		}
		for i := 0; i < 100; i++ {
			v.Erase(e, text.Region{A: 0, B: 1})
		}
		v.EndEdit(e)
	}
}

func BenchmarkViewScopeNameLinear(b *testing.B) {
	w := backend.GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	const (
		in     = "language_test.go"
		syntax = "testdata/package/Go.tmLanguage"
	)
	b.StopTimer()
	syn, err := newSyntax(syntax)
	if err != nil {
		b.Fatal(err)
	}
	backend.GetEditor().AddSyntax(syntax, syn)
	v.Settings().Set("syntax", syntax)
	if d, err := ioutil.ReadFile(in); err != nil {
		b.Fatal(err)
	} else {
		e := v.BeginEdit()
		v.Insert(e, 0, string(d))
		v.EndEdit(e)
		b.StartTimer()
		for j := 0; j < b.N; j++ {
			for i := 0; i < v.Size(); i++ {
				v.ScopeName(i)
			}
		}
	}
}

func BenchmarkViewScopeNameRandom(b *testing.B) {
	w := backend.GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	const (
		in     = "language_test.go"
		syntax = "testdata/package/Go.tmLanguage"
	)
	b.StopTimer()
	syn, err := newSyntax(syntax)
	if err != nil {
		b.Fatal(err)
	}
	backend.GetEditor().AddSyntax(syntax, syn)
	v.Settings().Set("syntax", syntax)
	if d, err := ioutil.ReadFile(in); err != nil {
		b.Fatal(err)
	} else {
		e := v.BeginEdit()
		v.Insert(e, 0, string(d))
		v.EndEdit(e)
		p := rand.Perm(b.N)
		b.StartTimer()
		for _, i := range p {
			v.ScopeName(i)
		}
	}
}
