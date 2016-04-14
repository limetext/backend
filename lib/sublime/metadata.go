// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/limetext/lime-backend/lib/loaders"
)

type (
	Metadata struct {
		Name     string
		Scope    string
		Settings MetaSettings `json:"settings"`
		UUID     string
	}

	MetaSettings struct {
		IncreaseIndentPattern        Regex
		DecreaseIndentPattern        Regex
		BracketIndentNextLinePattern Regex
		DisableIndentNextLinePattern Regex
		UnIndentedLinePattern        Regex
		CancelCompletion             Regex
		ShowInSymbolList             int
		ShowInIndexedSymbolList      int
		SymbolTransformation         Regex
		SymbolIndexTransformation    Regex
		ShellVariables               ShellVariables
	}

	ShellVariables map[string]string
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

func (m Metadata) String() string {
	ret := fmt.Sprintf("%s - %s\n", m.Name, m.UUID)
	ret += fmt.Sprintf("Scope: %s\n", m.Scope)
	ret += fmt.Sprintf("Settings\n%s", m.Settings)
	return ret
}

func (m MetaSettings) String() string {
	ret := fmt.Sprintf("\tIncreaseIndentPattern\n\t\t%s\n", m.IncreaseIndentPattern)
	ret += fmt.Sprintf("\tDecreaseIndentPattern\n\t\t%s\n", m.DecreaseIndentPattern)
	ret += fmt.Sprintf("\tBracketIndentNextLinePattern\n\t\t%s\n", m.BracketIndentNextLinePattern)
	ret += fmt.Sprintf("\tDisableIndentNextLinePattern\n\t\t%s\n", m.DisableIndentNextLinePattern)
	ret += fmt.Sprintf("\tUnIndentedLinePattern\n\t\t%s\n", m.UnIndentedLinePattern)
	ret += fmt.Sprintf("\tCancelCompletion\n\t\t%s\n", m.CancelCompletion)
	ret += fmt.Sprintf("\tShowInSymbolList\n\t\t%d\n", m.ShowInSymbolList)
	ret += fmt.Sprintf("\tShowInIndexedSymbolList\n\t\t%d\n", m.ShowInIndexedSymbolList)
	ret += fmt.Sprintf("\tSymbolTransformation\n\t\t%s\n", m.SymbolTransformation)
	ret += fmt.Sprintf("\tSymbolIndexTransformation\n\t\t%s\n", m.SymbolIndexTransformation)
	ret += fmt.Sprintf("\tShellVariables\n%s", m.ShellVariables)
	return ret
}

func (s *ShellVariables) UnmarshalJSON(data []byte) error {
	*s = make(ShellVariables)
	tmp := []struct {
		Name  string
		Value string
	}{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	for _, st := range tmp {
		(*s)[st.Name] = st.Value
	}
	return nil
}

func (s ShellVariables) String() (ret string) {
	keys := make([]string, 0, len(s))
	for k, _ := range s {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		ret += fmt.Sprintf("\t\t%s: '%s'\n", key, s[key])
	}
	return
}
