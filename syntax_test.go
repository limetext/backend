// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"testing"

	"github.com/limetext/backend/parser"
	"github.com/limetext/sublime/textmate/language"
	"github.com/limetext/text"
)

type dummySyntax struct {
	l *language.Language
}

func newDummySytax(tb testing.TB, path string) *dummySyntax {
	if l, err := language.Load(path); err != nil {
		tb.Fatalf("Error on loading language %s: %s", path, err)
		return nil
	} else {
		return &dummySyntax{l: l}
	}
}

func (s *dummySyntax) Parser(data string) (parser.Parser, error) {
	l := s.l.Copy()
	return language.NewParser(l, []rune(data)), nil
}

func (s *dummySyntax) Name() string {
	return s.l.Name
}

func (s *dummySyntax) FileTypes() []string {
	return s.l.FileTypes
}

func addSetSyntax(tb testing.TB, settings *text.Settings, path string) {
	syn := newDummySytax(tb, path)
	GetEditor().AddSyntax(path, syn)
	settings.Set("syntax", path)
}
