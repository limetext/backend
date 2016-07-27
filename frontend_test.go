// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"sync"

	"github.com/limetext/backend/log"
	"github.com/limetext/text"
)

type dummyFrontend struct {
	m sync.Mutex
	// Default return value for OkCancelDialog
	defaultAction bool
}

func (fe *dummyFrontend) SetDefaultAction(action bool) {
	fe.m.Lock()
	defer fe.m.Unlock()
	fe.defaultAction = action
}
func (fe *dummyFrontend) StatusMessage(msg string) { log.Info(msg) }
func (fe *dummyFrontend) ErrorMessage(msg string)  { log.Error(msg) }
func (fe *dummyFrontend) MessageDialog(msg string) { log.Info(msg) }
func (fe *dummyFrontend) OkCancelDialog(msg string, button string) bool {
	log.Info(msg)
	fe.m.Lock()
	defer fe.m.Unlock()
	return fe.defaultAction
}
func (fe *dummyFrontend) Show(v *View, r text.Region) {}
func (fe *dummyFrontend) VisibleRegion(v *View) text.Region {
	return text.Region{}
}
func (fe *dummyFrontend) Prompt(title, folder string, flags int) []string {
	return nil
}
