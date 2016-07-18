// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/limetext/text"
)

type (
	Project struct {
		text.HasSettings
		filename string
		folders  Folders
		// TODO: build_systems
	}

	Folder struct {
		Path                string   `json:"path"`
		Name                string   `json:"name"`
		ExcludePatterns     []string `json:"folder_exclude_patterns"`
		IncludePatterns     []string `json:"folder_include_patterns"`
		FileExcludePatterns []string `json:"file_exclude_patterns"`
		FileIncludePatterns []string `json:"file_include_patterns"`
		FollowSymlinks      bool     `json:"follow_symlinks"`
	}

	Folders []*Folder
)

func New() *Project {
	p := &Project{}
	p.folders = make(Folders, 0)
	return p
}

func (p *Project) Close() {
	p = New()
}

func (p *Project) SaveAs(name string) error {
	if data, err := json.Marshal(p); err != nil {
		return err
	} else if err := ioutil.WriteFile(name, data, 0644); err != nil {
		return err
	}
	if abs, err := filepath.Abs(name); err != nil {
		p.SetName(name)
	} else {
		p.SetName(abs)
	}
	return nil
}

func (p *Project) AddFolder(name string) {
	p.folders = append(p.folders, &Folder{Path: name})
}

func (p *Project) RemoveFolder(name string) {
	for i, folder := range p.folders {
		if name == folder.Path {
			p.folders = append(p.folders[:i], p.folders[i+1:]...)
			break
		}
	}
}

func (p *Project) Folders() []string {
	folders := make([]string, 0, len(p.folders))
	for _, folder := range p.folders {
		folders = append(folders, folder.Path)
	}
	return folders
}

func (p *Project) Folder(path string) *Folder {
	for _, folder := range p.folders {
		if folder.Path == path {
			return folder
		}
	}
	return nil
}

func (p *Project) FileName() string {
	return p.filename
}

func (p *Project) SetName(name string) {
	p.filename = name
}

func (p *Project) UnmarshalJSON(data []byte) error {
	med := struct {
		Folders  Folders
		Settings text.Settings
	}{}
	if err := json.Unmarshal(data, &med); err != nil {
		return err
	}
	p.folders = med.Folders
	if data, err := json.Marshal(&med.Settings); err != nil {
		return err
	} else if err = json.Unmarshal(data, p.Settings()); err != nil {
		return err
	}
	return nil
}

func (p *Project) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("{\n\t\"folders\":\n\t[\n")
	for i, folder := range p.folders {
		if i != 0 {
			buf.WriteString(",\n")
		}
		buf.WriteString("\t\t{\n")
		fmt.Fprintf(buf, "\t\t\t\"path\": \"%s\"", folder.Path)
		if folder.Name != "" {
			fmt.Fprintf(buf, ",\n\t\t\t\"name\": \"%s\"", folder.Name)
		}
		if folder.ExcludePatterns != nil {
			fmt.Fprintf(buf, ",\n\t\t\t\"folder_exclude_patterns\": \"%s\"", folder.ExcludePatterns)
		}
		if folder.IncludePatterns != nil {
			fmt.Fprintf(buf, ",\n\t\t\t\"folder_include_patterns\": \"%s\"", folder.IncludePatterns)
		}
		if folder.FileExcludePatterns != nil {
			fmt.Fprintf(buf, ",\n\t\t\t\"file_exclude_patterns\": \"%s\"", folder.ExcludePatterns)
		}
		if folder.FileIncludePatterns != nil {
			fmt.Fprintf(buf, ",\n\t\t\t\"file_include_patterns\": \"%s\"", folder.ExcludePatterns)
		}
		if folder.FollowSymlinks {
			fmt.Fprintf(buf, ",\n\t\t\t\"follow_symlinks\": \"%t\"", folder.FollowSymlinks)
		}
		buf.WriteString("\n\t\t}")
	}
	buf.WriteString("\n\t]")
	if data, err := json.MarshalIndent(p.Settings(), "", "\t"); err != nil {
		return nil, err
	} else if str := string(data); str != "{}" {
		str = strings.Replace(str, "\t", "\t\t", -1)
		str = strings.Replace(str, "{", "\t{", -1)
		str = strings.Replace(str, "}", "\t}", -1)
		fmt.Fprintf(buf, ",\n\t\"settings\":\n%s", str)
	}
	buf.WriteString("\n}\n")
	return buf.Bytes(), nil
}
