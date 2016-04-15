// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	year          = strconv.FormatInt(int64(time.Now().Year()), 10)
	licenseheader = []byte(`// Copyright ` + year + ` The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.
`)
)

func license(path string, fi os.FileInfo, err error) error {
	if fi.IsDir() {
		switch filepath.Base(path) {
		case "testdata":
			return filepath.SkipDir
		case "vendor":
			return filepath.SkipDir
		}
		return nil
	}

	switch filepath.Ext(path) {
	case ".go":
	default:
		return nil
	}

	changed := false
	cmp, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	lhn := append(licenseheader, '\n')
	if !bytes.Equal([]byte("// Copyright"), cmp[:12]) {
		cmp = append(lhn, cmp...)
		log.Println("Added license to", path)
		changed = true
	}

	if changed {
		return ioutil.WriteFile(path, cmp, fi.Mode().Perm())
	}

	return nil
}

var path = flag.String("path", "./", "path to run fix command")

func main() {
	if err := filepath.Walk(*path, license); err != nil {
		log.Println("done: ", err)
	}
}
