// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"fmt"
	"io/ioutil"

	"github.com/limetext/lime-backend/lib/loaders"
)

type (
	Metadata struct {
		Name     string
		Scope    string
		Settings Settings
		UUID     string
	}

	Settings struct {
		// string is regex
		IncreaseIndentPattern        string
		DecreaseIndentPattern        string
		BracketIndentNextLinePattern string
		DisableIndentNextLinePattern string
		UnIndentedLinePattern        string
		CancelCompletion             string
		ShowInSymbolList             int
		ShowInIndexedSymbolList      int
		SymbolTransformation         string
		SymbolIndexTransformation    string
		ShellVariables               ShellVariable
	}

	ShellVariable map[string]string
)

func LoadMetadata(filename string) (*Metadata, error) {
	var md Metadata
	if d, err := ioutil.ReadFile(filename); err != nil {
		return nil, fmt.Errorf("Unable to read metadata file: %s", err)
	} else if err = loaders.LoadPlist(d, &md); err != nil {
		return nil, fmt.Errorf("Unable to load metadata data: %s", err)
	}

	return &md, nil
}

func (m *Metadata) String() string {
	ret := fmt.Sprintln("%s - %s", m.Name, m.UUID)
	ret += fmt.Sprintln("\tScope\n\t\t%s", m.Scope)
	ret += fmt.Sprintf("%s", m.Settings)

	return ret
}
