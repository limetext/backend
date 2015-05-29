// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package textmate

import (
	"fmt"
	"io/ioutil"

	"github.com/limetext/lime-backend/lib/loaders"
)

type (
	// http://sublime-text-unofficial-documentation.readthedocs.org/en/latest/reference/metadata.html
	Preferences struct {
		BundleUUID UUID
		Name       string
		Scope      string
		Settings   PrefSetting
		UUID       UUID
	}

	PrefSetting struct {
		ShellVariables               []ShellVariable
		SmartTypingPairs             Pairs
		HighlightPairs               Pairs
		IncreaseIndentPattern        Regex
		DecreaseIndentPattern        Regex
		UnIndentedLinePattern        Regex
		BracketIndentNextLinePattern Regex
		IndentNextLinePattern        Regex
		DisableIndentNextLinePattern Regex
		ZeroIndentPattern            Regex
		DisableIndentCorrections     string
		FoldingStartMarker           Regex
		FoldingStopMarker            Regex
		FoldingIndentedBlockStart    Regex
		FoldingIndentedBlockIgnore   Regex
		Completions                  []string
		CompletionCommand            string
		CancelCompletion             Regex
		DisableDefaultCompletion     int
		ShowInSymbolList             int
		ShowInIndexedSymbolList      int
		SymbolTransformation         Regex
		SymbolIndexTransformation    Regex
		SpellChecking                bool
	}

	ShellVariable struct {
		Name  string
		Value string
	}

	Pairs []Pair

	Pair [2]string
)

func LoadPrefrences(filename string) (*Preferences, error) {
	var pref Preferences
	if d, err := ioutil.ReadFile(filename); err != nil {
		return nil, fmt.Errorf("Unable to load preference definition: %s", err)
	} else if err := loaders.LoadPlist(d, &pref); err != nil {
		return nil, fmt.Errorf("Unable to load preference definition: %s", err)
	}

	return &pref, nil
}
