// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"time"

	"github.com/limetext/gopy/lib"
	"github.com/limetext/lime-backend/lib"
	"github.com/limetext/lime-backend/lib/log"
	"github.com/limetext/lime-backend/lib/render"
	"github.com/limetext/lime-backend/lib/util"
)

func sublime_Console(tu *py.Tuple, kwargs *py.Dict) (py.Object, error) {
	if tu.Size() != 1 {
		return nil, fmt.Errorf("Unexpected argument count: %d", tu.Size())
	}
	if i, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		log.Info("Python sez: %s", i)
	}
	return toPython(nil)
}

func sublime_set_timeout(tu *py.Tuple, kwargs *py.Dict) (py.Object, error) {
	var (
		pyarg py.Object
	)
	if tu.Size() != 2 {
		return nil, fmt.Errorf("Unexpected argument count: %d", tu.Size())
	}
	if i, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		pyarg = i
	}
	if i, err := tu.GetItem(1); err != nil {
		return nil, err
	} else if v, err := fromPython(i); err != nil {
		return nil, err
	} else if v2, ok := v.(int); !ok {
		return nil, fmt.Errorf("Expected int not %s", i.Type())
	} else {
		pyarg.Incref()
		go func() {
			time.Sleep(time.Millisecond * time.Duration(v2))
			l := py.NewLock()
			defer l.Unlock()
			defer pyarg.Decref()
			if ret, err := pyarg.Base().CallFunctionObjArgs(); err != nil {
				log.Debug("Error in callback: %v", err)
			} else {
				ret.Decref()
			}
		}()
	}
	return toPython(nil)
}

func sublime_PackagesPath(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 string
	)
	if tu.Size() == 0 {
		arg1 = "shipped"
	} else if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(string); !ok {
				return nil, fmt.Errorf("Expected type string for backend.Editor.PackagesPath() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	ret0 := backend.GetEditor().PackagesPath(arg1)
	var err error
	var pyret0 py.Object

	pyret0, err = toPython(ret0)
	if err != nil {
		return nil, err
	}
	return pyret0, err
}

var sublime_manual_methods = []py.Method{
	{Name: "console", Func: sublime_Console},
	{Name: "set_timeout", Func: sublime_set_timeout},
	{Name: "packages_path", Func: sublime_PackagesPath},
}

// TODO: check how many times is this function running
func init() {
	sublime_methods = append(sublime_methods, sublime_manual_methods...)
	l := py.InitAndLock()
	defer l.Unlock()

	m, err := py.InitModule("sublime", sublime_methods)
	if err != nil {
		// TODO: we should handle this as error
		panic(err)
	}

	if sys, err := py.Import("sys"); err != nil {
		log.Warn(err)
	} else {
		if pyc, err := py.NewUnicode("dont_write_bytecode"); err != nil {
			log.Warn(err)
		} else {
			// avoid pyc files
			sys.Base().SetAttr(pyc, py.True)
		}
		sys.Decref()
	}

	classes := []struct {
		name string
		c    *py.Class
	}{
		{"Region", &_regionClass},
		{"RegionSet", &_region_setClass},
		{"View", &_viewClass},
		{"Window", &_windowClass},
		{"Edit", &_editClass},
		{"Settings", &_settingsClass},
		{"WindowCommandGlue", &_windowCommandGlueClass},
		{"TextCommandGlue", &_textCommandGlueClass},
		{"ApplicationCommandGlue", &_applicationCommandGlueClass},
		{"OnQueryContextGlue", &_onQueryContextGlueClass},
		{"ViewEventGlue", &_viewEventGlueClass},
	}
	constants := []struct {
		name     string
		constant int
	}{
		{"OP_EQUAL", int(util.OpEqual)},
		{"OP_NOT_EQUAL", int(util.OpNotEqual)},
		{"OP_REGEX_MATCH", int(util.OpRegexMatch)},
		{"OP_NOT_REGEX_MATCH", int(util.OpNotRegexMatch)},
		{"OP_REGEX_CONTAINS", int(util.OpRegexContains)},
		{"OP_NOT_REGEX_CONTAINS", int(util.OpNotRegexContains)},
		{"INHIBIT_WORD_COMPLETIONS", 0},
		{"INHIBIT_EXPLICIT_COMPLETIONS", 0},
		{"LITERAL", int(backend.IGNORECASE)},
		{"IGNORECASE", int(backend.LITERAL)},
		{"CLASS_WORD_START", int(backend.CLASS_WORD_START)},
		{"CLASS_WORD_END", int(backend.CLASS_WORD_END)},
		{"CLASS_PUNCTUATION_START", int(backend.CLASS_PUNCTUATION_START)},
		{"CLASS_PUNCTUATION_END", int(backend.CLASS_PUNCTUATION_END)},
		{"CLASS_SUB_WORD_START", int(backend.CLASS_SUB_WORD_START)},
		{"CLASS_SUB_WORD_END", int(backend.CLASS_SUB_WORD_END)},
		{"CLASS_LINE_START", int(backend.CLASS_LINE_START)},
		{"CLASS_LINE_END", int(backend.CLASS_LINE_END)},
		{"CLASS_EMPTY_LINE", int(backend.CLASS_EMPTY_LINE)},
		{"CLASS_MIDDLE_WORD", int(backend.CLASS_MIDDLE_WORD)},
		{"CLASS_WORD_START_WITH_PUNCTUATION", int(backend.CLASS_WORD_START_WITH_PUNCTUATION)},
		{"CLASS_WORD_END_WITH_PUNCTUATION", int(backend.CLASS_WORD_END_WITH_PUNCTUATION)},
		{"CLASS_OPENING_PARENTHESIS", int(backend.CLASS_OPENING_PARENTHESIS)},
		{"CLASS_CLOSING_PARENTHESIS", int(backend.CLASS_CLOSING_PARENTHESIS)},
		{"DRAW_EMPTY", int(render.DRAW_EMPTY)},
		{"HIDE_ON_MINIMAP", int(render.HIDE_ON_MINIMAP)},
		{"DRAW_EMPTY_AS_OVERWRITE", int(render.DRAW_EMPTY_AS_OVERWRITE)},
		{"DRAW_NO_FILL", int(render.DRAW_NO_FILL)},
		{"DRAW_NO_OUTLINE", int(render.DRAW_NO_OUTLINE)},
		{"DRAW_SOLID_UNDERLINE", int(render.DRAW_SOLID_UNDERLINE)},
		{"DRAW_STIPPLED_UNDERLINE", int(render.DRAW_STIPPLED_UNDERLINE)},
		{"DRAW_SQUIGGLY_UNDERLINE", int(render.DRAW_SQUIGGLY_UNDERLINE)},
		{"PERSISTENT", int(render.PERSISTENT)},
		{"HIDDEN", int(render.HIDDEN)},
	}

	for _, cl := range classes {
		c, err := cl.c.Create()
		if err != nil {
			panic(err)
		}
		if err := m.AddObject(cl.name, c); err != nil {
			panic(err)
		}
	}
	for _, c := range constants {
		if err := m.AddIntConstant(c.name, c.constant); err != nil {
			panic(err)
		}
	}
}
