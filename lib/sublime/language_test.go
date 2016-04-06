// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	"github.com/limetext/lime-backend/lib"
	"github.com/limetext/lime-backend/lib/util"
	"github.com/limetext/text"
)

func TestLanguageProviderLanguageFromScope(t *testing.T) {
	l, _ := Provider.LanguageFromFile("testdata/Go.tmLanguage")

	if _, err := Provider.LanguageFromScope(l.ScopeName); err != nil {
		t.Errorf("Tried to load %s, but got an error: %v", l.ScopeName, err)
	}

	if _, err := Provider.LanguageFromScope("MissingScope"); err == nil {
		t.Error("Tried to load MissingScope, expecting to get an error, but didn't")
	}
}

func TestLanguageProviderLanguageFromFile(t *testing.T) {
	if _, err := Provider.LanguageFromFile("testdata/Go.tmLanguage"); err != nil {
		t.Errorf("Tried to load testdata/Go.tmLanguage, but got an error: %v", err)
	}

	if _, err := Provider.LanguageFromFile("MissingFile"); err == nil {
		t.Error("Tried to load MissingFile, expecting to get an error, but didn't")
	}
}

func TestTmLanguage(t *testing.T) {
	files := []string{
		"testdata/Property List (XML).tmLanguage",
		"testdata/XML.plist",
		"testdata/Go.tmLanguage",
	}
	for _, fn := range files {
		if _, err := Provider.LanguageFromFile(fn); err != nil {
			t.Fatal(err)
		}
	}

	type test struct {
		in  string
		out string
		syn string
	}
	tests := []test{
		{
			"testdata/plist.tmlang",
			"testdata/plist.tmlang.res",
			"text.xml.plist",
		},
		{
			"testdata/Property List (XML).tmLanguage",
			"testdata/Property List (XML).tmLanguage.res",
			"text.xml.plist",
		},
		{
			"testdata/main.go",
			"testdata/main.go.res",
			"source.go",
		},
		{
			"testdata/go2.go",
			"testdata/go2.go.res",
			"source.go",
		},
		{
			"testdata/utf.go",
			"testdata/utf.go.res",
			"source.go",
		},
	}
	for _, t3 := range tests {

		var d0 string
		if d, err := ioutil.ReadFile(t3.in); err != nil {
			t.Errorf("Couldn't load file %s: %s", t3.in, err)
			continue
		} else {
			d0 = string(d)
		}

		if syn, err := syntaxFromLanguage(t3.syn); err != nil {
			t.Error(err)
		} else if pr, err := syn.Parser(d0); err != nil {
			t.Error(err)
		} else if root, err := pr.Parse(); err != nil {
			t.Error(err)
		} else {
			str := fmt.Sprintf("%s", root)
			if d, err := ioutil.ReadFile(t3.out); err != nil {
				if err := ioutil.WriteFile(t3.out, []byte(str), 0644); err != nil {
					t.Error(err)
				}
			} else if diff := util.Diff(string(d), str); diff != "" {
				t.Error(diff)
			}
		}
	}
}

func BenchmarkLanguage(b *testing.B) {
	b.StopTimer()
	tst := []string{
		"language.go",
		"testdata/main.go",
	}

	var d0 []string
	for _, s := range tst {
		if d, err := ioutil.ReadFile(s); err != nil {
			b.Errorf("Couldn't load file %s: %s", s, err)
		} else {
			d0 = append(d0, string(d))
		}
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := range d0 {
			syn, err := newSyntax("testdata/Go.tmLanguage")
			if err != nil {
				b.Fatal(err)
				return
			}
			pr, err := syn.Parser(d0[j])
			if err != nil {
				b.Fatal(err)
				return
			}
			_, err = pr.Parse()
			if err != nil {
				b.Fatal(err)
				return
			}
		}
	}
	fmt.Println(util.Prof)
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
		syntax  = "testdata/Go.tmLanguage"
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
		syntax  = "testdata/Go.tmLanguage"
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

	syntax := "testdata/Go.tmLanguage"
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
		syntax = "testdata/Go.tmLanguage"
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
		syntax = "testdata/Go.tmLanguage"
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

func syntaxFromLanguage(id string) (*Syntax, error) {
	l, err := Provider.GetLanguage(id)
	if err != nil {
		return nil, err
	}
	return &Syntax{l: l}, nil
}
