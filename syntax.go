// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"fmt"

	"github.com/limetext/backend/log"
	"github.com/limetext/backend/parser"
	"github.com/limetext/backend/render"
	"github.com/limetext/text"
)

// Any syntax definition for view should implement this interface
// also it should register it self from editor.AddSyntax
type Syntax interface {
	// provides parser for creating syntax highlighter
	Parser(data string) (parser.Parser, error)
	Name() string
	// filetypes this syntax supports
	FileTypes() []string
}

func syntaxHighlighter(name, data string) parser.SyntaxHighlighter {
	if name == "" {
		return &syntax{}
	}
	sh, err := syntaxProvider(name, data)
	if err != nil {
		log.Error("%s, falling back to default syntax", err)
		return &syntax{}
	}
	return sh
}

func syntaxProvider(name, data string) (parser.SyntaxHighlighter, error) {
	syn := GetEditor().GetSyntax(name)
	if syn == nil {
		return nil, fmt.Errorf("No syntax %s in editor", name)
	}
	pr, err := syn.Parser(data)
	if err != nil {
		return nil, fmt.Errorf("Couldn't get parser from syntax: %s", err)
	}
	sh, err := parser.NewSyntaxHighlighter(pr)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create syntaxhighlighter: %s", err)
	}
	return sh, nil
}

type syntax struct{}

func (s *syntax) Adjust(position, delta int) {}

func (s *syntax) ScopeExtent(point int) text.Region {
	return text.Region{}
}

func (s *syntax) ScopeName(p int) string {
	return "text.plain"
}

func (s *syntax) Flatten() render.ViewRegionMap {
	return nil
}
