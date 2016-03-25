// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package textmate

type (
	// http://docs.sublimetext.info/en/latest/reference/metadata.html
	Preferences struct {
		Name     string
		Scope    string
		Settings PrefSetting
		UUID     UUID
	}

	PrefSetting struct {
		ShellVariables []ShellVariable
	}

	ShellVariable struct {
		Name  string
		Value string
	}
)
