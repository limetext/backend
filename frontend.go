// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import "github.com/limetext/text"

// The Frontend interface defines the API
// for functionality that is frontend specific.
type Frontend interface {
	// Probe the frontend for the currently
	// visible region of the given view.
	VisibleRegion(v *View) text.Region

	// Make the frontend show the specified region of the
	// given view.
	Show(v *View, r text.Region)

	// Sets the status message shown in the status bar
	StatusMessage(string)

	// Displays an error message to the user
	ErrorMessage(string)

	// Displays a message dialog to the user
	MessageDialog(string)

	// Displays an ok / cancel dialog to the user.
	// "okname" if provided will be used as the text
	// instead of "Ok" for the ok button.
	// Returns true when ok was pressed, and false when
	// cancel was pressed.
	OkCancelDialog(msg string, okname string) bool

	// Displays file dialog, returns the selected files.
	// folder is the path file dialog will show.
	Prompt(title, folder string, flags int) []string
}

const (
	// Prompt save as dialog
	SaveAs = 1 << iota
	// User should only be able to select folders
	OnlyFolder
	// User can select multiple files
	SelectMultiple
)
