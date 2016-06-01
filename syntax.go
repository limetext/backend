// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"fmt"

	"github.com/limetext/backend/parser"
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

func syntaxHighlighter(name, data string) (parser.SyntaxHighlighter, error) {
	if name == "" {
		// TODO bring the default syntax up
	}
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
		return nil, fmt.Errorf("Couldn't create syntaxhighlighter: %v", err)
	}
	return sh, nil
}
