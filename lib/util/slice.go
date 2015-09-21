// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package util

// Checks if element exists in a slice of strings
func Exists(strs []string, s string) bool {
	for _, str := range strs {
		if str == s {
			return true
		}
	}
	return false
}

// Removes an element from slice of strings
func Remove(strs []string, s string) []string {
	for i, el := range strs {
		if el == s {
			strs[i], strs = strs[len(strs)-1], strs[:len(strs)-1]
			break
		}
	}
	return strs
}
