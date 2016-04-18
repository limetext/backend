// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

// The packages package handles lime package management.
//
// The key idea of lime packages is modularity/optionality. The core backend
// shouldn't know about tmbundle nor sublime-package, but rather it should
// make it possible to use these and other variants. @quarnster
// Ideally packages implemented in such a way that we can just do:
// import (
// _ "github.com/limetext/lime/backend/textmate"
// _ "github.com/limetext/lime/backend/sublime"
// _ "github.com/limetext/lime/backend/emacs"
// )
//
// Package type
//
// Each plugin or package that wants to communicate with backend should
// implement this interface.
//
// Record type
//
// For enabling lime to detect and load a package it should register itself as
// a Record
//
package packages
