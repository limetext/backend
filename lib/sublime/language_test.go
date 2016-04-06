// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/limetext/lime-backend/lib/util"
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

func syntaxFromLanguage(id string) (*Syntax, error) {
	l, err := Provider.GetLanguage(id)
	if err != nil {
		return nil, err
	}
	return &Syntax{l: l}, nil
}
