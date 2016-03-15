// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/limetext/gopy/lib"
	"github.com/limetext/lime-backend/lib/util"
)

// TODO: this could be much better
// "__*__" should be omitted and it should be same as
// https://www.sublimetext.com/docs/3/api_reference.html
func TestSublimeApi(t *testing.T) {
	const expfile = "testdata/api.txt"
	l := py.NewLock()
	defer l.Unlock()
	subl, err := py.Import("sublime")
	if err != nil {
		t.Fatal(err)
	}
	buf := bytes.NewBuffer(nil)

	if err := printObj("", subl, buf); err != nil {
		t.Error(err)
	}
	if exp, err := ioutil.ReadFile(expfile); err != nil {
		t.Fatal(err)
	} else if diff := util.Diff(string(exp), buf.String()); diff != "" {
		t.Error(diff)
	}
}

func printObj(indent string, v py.Object, buf *bytes.Buffer) error {
	b := v.Base()
	dir, err := b.Dir()
	if err != nil {
		return err
	}
	l, ok := dir.(*py.List)
	if !ok {
		return fmt.Errorf("Unexpected type: %v", dir.Type())
	}
	sl := l.Slice()
	if indent == "" {
		for _, v2 := range sl {
			if item, err := b.GetAttr(v2); err != nil {
				return err
			} else {
				ty := item.Type()
				line := fmt.Sprintf("%s%s\n", indent, v2)
				buf.WriteString(line)
				if ty == py.TypeType {
					if err := printObj(indent+"\t", item, buf); err != nil {
						return err
					}
				}
				item.Decref()
			}
		}
	} else {
		for _, v2 := range sl {
			buf.WriteString(fmt.Sprintf("%s%s\n", indent, v2))
		}
	}
	dir.Decref()
	return nil
}
