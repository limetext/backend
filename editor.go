// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"fmt"
	"path"
	"runtime"
	"runtime/debug"
	"sync"

	"github.com/limetext/backend/clipboard"
	"github.com/limetext/backend/keys"
	"github.com/limetext/backend/log"
	"github.com/limetext/backend/packages"
	"github.com/limetext/backend/watch"
	"github.com/limetext/text"
	"github.com/limetext/util"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	GetEditor()
}

type Editor struct {
	text.HasSettings
	keys.HasKeyBindings
	*watch.Watcher
	windows          []*Window
	activeWindow     *Window
	logInput         bool
	cmdHandler       commandHandler
	console          *View
	frontend         Frontend
	keyInput         chan (keys.KeyPress)
	clipboard        clipboard.Clipboard
	defaultSettings  *text.HasSettings
	platformSettings *text.HasSettings
	defaultKB        *keys.HasKeyBindings
	platformKB       *keys.HasKeyBindings
	userKB           *keys.HasKeyBindings
	defaultPath      string
	userPath         string
	pkgsPaths        []string
	colorSchemes     map[string]ColorScheme
	syntaxes         map[string]Syntax
	filetypes        map[string]string
}

var (
	ed  *Editor
	edl sync.Mutex
)

func GetEditor() *Editor {
	edl.Lock()
	defer edl.Unlock()
	if ed == nil {
		ed = &Editor{
			cmdHandler: commandHandler{
				ApplicationCommands: make(appcmd),
				TextCommands:        make(textcmd),
				WindowCommands:      make(wndcmd),
				verbose:             true,
			},
			console: &View{
				buffer:  text.NewBuffer(),
				scratch: true,
			},
			keyInput:         make(chan keys.KeyPress, 32),
			clipboard:        clipboard.NewSystemClipboard(),
			defaultSettings:  new(text.HasSettings),
			platformSettings: new(text.HasSettings),
			defaultKB:        new(keys.HasKeyBindings),
			platformKB:       new(keys.HasKeyBindings),
			userKB:           new(keys.HasKeyBindings),
			pkgsPaths:        make([]string, 0),
			colorSchemes:     make(map[string]ColorScheme),
			syntaxes:         make(map[string]Syntax),
			filetypes:        make(map[string]string),
		}
		var err error
		if ed.Watcher, err = watch.NewWatcher(); err != nil {
			log.Error("Couldn't create watcher: %s", err)
		}

		ed.console.Settings().Set("is_widget", true)
		// Initializing settings hierarchy
		// default <- platform <- user(editor)
		ed.platformSettings.Settings().SetParent(ed.defaultSettings)
		ed.Settings().SetParent(ed.platformSettings)

		// Initializing keybidings hierarchy
		// default <- platform <- user <- user platform(editor)
		ed.KeyBindings().SetParent(ed.userKB)
		ed.userKB.KeyBindings().SetParent(ed.platformKB)
		ed.platformKB.KeyBindings().SetParent(ed.defaultKB)

		OnDefaultPathAdd.Add(ed.loadDefaultSettings)
		OnDefaultPathAdd.Add(ed.loadDefaultKeyBindings)
		OnUserPathAdd.Add(ed.loadUserSettings)
		OnUserPathAdd.Add(ed.loadUserKeyBindings)
		ed.Settings().AddOnChange("backend.editor.ignored_packages", func(name string) {
			if name != "ignored_packages" {
				return
			}
			// If a package is removed from ignored_packages it will
			// be loaded by re scanning
			for _, path := range ed.pkgsPaths {
				packages.Scan(path)
			}

			ignoreds, ok := ed.Settings().Get("ignored_packages").([]interface{})
			if !ok {
				return
			}
			for _, ignored := range ignoreds {
				if name, ok := ignored.(string); ok {
					packages.UnLoad(name)
				}
			}
		})

		log.AddFilter("console", log.DEBUG, log.NewLogWriter(ed.handleLog))
		go ed.inputthread()
	}
	return ed
}

func (e *Editor) Frontend() Frontend {
	return e.frontend
}

func (e *Editor) SetFrontend(f Frontend) {
	e.frontend = f
}

func (e *Editor) Init() {
	log.Info("Initializing")
	OnInit.call()
	OnPackagesPathAdd.Add(packages.Scan)
}

func (e *Editor) loadDefaultKeyBindings(dir string) {
	log.Fine("Loading editor default keybindings")
	p := path.Join(dir, "Default.sublime-keymap")
	log.Finest("Loading %s", p)
	packages.LoadJSON(p, e.defaultKB.KeyBindings())

	p = path.Join(dir, "Default ("+e.Plat()+").sublime-keymap")
	log.Finest("Loading %s", p)
	packages.LoadJSON(p, e.platformKB.KeyBindings())
}

func (e *Editor) loadUserKeyBindings(dir string) {
	log.Fine("Loading editor user keybindings")
	p := path.Join(dir, "Default.sublime-keymap")
	log.Finest("Loading %s", p)
	packages.LoadJSON(p, e.userKB.KeyBindings())

	p = path.Join(dir, "Default ("+e.Plat()+").sublime-keymap")
	log.Finest("Loading %s", p)
	packages.LoadJSON(p, e.KeyBindings())
}

func (e *Editor) loadDefaultSettings(dir string) {
	log.Fine("Loading editor default settings")
	p := path.Join(dir, "Preferences.sublime-settings")
	log.Finest("Loading %s", p)
	packages.LoadJSON(p, e.defaultSettings.Settings())

	p = path.Join(dir, "Preferences ("+e.Plat()+").sublime-settings")
	log.Finest("Loading %s", p)
	packages.LoadJSON(p, e.platformSettings.Settings())
}

func (e *Editor) loadUserSettings(dir string) {
	log.Fine("Loading editor user settings")
	p := path.Join(dir, "Preferences.sublime-settings")
	log.Finest("Loading %s", p)
	packages.LoadJSON(p, e.Settings())
}

func (e *Editor) PackagesPath() string {
	if len(e.pkgsPaths) == 0 {
		return ""
	}
	return e.pkgsPaths[0]
}

func (e *Editor) Console() *View {
	return e.console
}

func (e *Editor) Windows() []*Window {
	edl.Lock()
	defer edl.Unlock()
	ret := make([]*Window, len(e.windows))
	copy(ret, e.windows)
	return ret
}

func (e *Editor) SetActiveWindow(w *Window) {
	e.activeWindow = w
}

func (e *Editor) ActiveWindow() *Window {
	return e.activeWindow
}

func (e *Editor) NewWindow() *Window {
	edl.Lock()
	e.windows = append(e.windows, &Window{})
	w := e.windows[len(e.windows)-1]
	edl.Unlock()
	w.Settings().SetParent(e)
	w.Project().Settings().SetParent(w)
	e.SetActiveWindow(w)
	OnNewWindow.Call(w)
	return w
}

func (e *Editor) remove(w *Window) {
	edl.Lock()
	defer edl.Unlock()
	for i, ww := range e.windows {
		if w == ww {
			end := len(e.windows) - 1
			if i != end {
				copy(e.windows[i:], e.windows[i+1:])
			}
			e.windows = e.windows[:end]
			return
		}
	}
	log.Error("Wanted to remove window %s, but it doesn't appear to be a child of this editor", w)
}

func (e *Editor) Arch() string {
	return runtime.GOARCH
}

func (e *Editor) Platform() string {
	return runtime.GOOS
}

func (e *Editor) Plat() string {
	switch e.Platform() {
	case "windows":
		return "Windows"
	case "darwin":
		return "OSX"
	}
	return "Linux"
}

func (e *Editor) Version() string {
	return "0"
}

func (e *Editor) CommandHandler() CommandHandler {
	return &e.cmdHandler
}

func (e *Editor) HandleInput(kp keys.KeyPress) {
	e.keyInput <- kp
}

func (e *Editor) inputthread() {
	pc := 0
	var lastBindings keys.KeyBindings
	doinput := func(kp keys.KeyPress) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Panic in inputthread: %v\n%s", r, string(debug.Stack()))
				if pc > 0 {
					panic(r)
				}
				pc++
			}
		}()
		p := util.Prof.Enter("hi")
		defer p.Exit()

		lvl := log.FINE
		if e.logInput {
			lvl++
		}
		log.Logf(lvl, "Key: %v", kp)
		if lastBindings.SeqIndex() == 0 {
			lastBindings = *e.KeyBindings()
		}
	try_again:
		possible_actions := lastBindings.Filter(kp)
		lastBindings = possible_actions

		// TODO?
		var (
			wnd *Window
			v   *View
		)
		if wnd = e.ActiveWindow(); wnd != nil {
			v = wnd.ActiveView()
		}

		qc := func(key string, operator util.Op, operand interface{}, match_all bool) bool {
			return OnQueryContext.Call(v, key, operator, operand, match_all) == True
		}

		if action := possible_actions.Action(qc); action != nil {
			p2 := util.Prof.Enter("hi.perform")
			e.RunCommand(action.Command, action.Args)
			p2.Exit()
		} else if possible_actions.SeqIndex() > 1 {
			// TODO: this disables having keyBindings with more than 2 key sequence
			lastBindings = *e.KeyBindings()
			goto try_again
		} else if kp.IsCharacter() {
			p2 := util.Prof.Enter("hi.character")
			log.Finest("[editor.inputthread] kp: |%s|, pos: %v", kp.Text, possible_actions)
			if err := e.CommandHandler().RunTextCommand(v, "insert", Args{"characters": kp.Text}); err != nil {
				log.Debug("Couldn't run textcommand: %s", err)
			}
			p2.Exit()
		}
	}
	for kp := range e.keyInput {
		doinput(kp)
	}
}

func (e *Editor) LogInput(l bool) {
	e.logInput = l
}

func (e *Editor) LogCommands(l bool) {
	e.cmdHandler.log = l
}

func (e *Editor) RunCommand(name string, args Args) {
	// TODO?
	var (
		wnd *Window
		v   *View
	)
	if wnd = e.ActiveWindow(); wnd != nil {
		v = wnd.ActiveView()
	}

	// TODO: what's the command precedence?
	if c := e.cmdHandler.TextCommands[name]; c != nil {
		if err := e.CommandHandler().RunTextCommand(v, name, args); err != nil {
			log.Debug("Couldn't run textcommand: %s", err)
		}
	} else if c := e.cmdHandler.WindowCommands[name]; c != nil {
		if err := e.CommandHandler().RunWindowCommand(wnd, name, args); err != nil {
			log.Debug("Couldn't run windowcommand: %s", err)
		}
	} else if c := e.cmdHandler.ApplicationCommands[name]; c != nil {
		if err := e.CommandHandler().RunApplicationCommand(name, args); err != nil {
			log.Debug("Couldn't run applicationcommand: %s", err)
		}
	} else {
		log.Debug("Couldn't find command to run")
	}
}

// Clipboard returns the currently active clipboard.
func (e *Editor) Clipboard() clipboard.Clipboard {
	return e.clipboard
}

// UseClipboard replaces the existing clipboard with c.
func (e *Editor) UseClipboard(c clipboard.Clipboard) {
	e.clipboard = c
}

// GetClipboard returns the contents of the clipboard. It assumes the text was
// not captured from an auto-expanded cursor. It exists for Sublime Text API
// compatibility.
func (e *Editor) GetClipboard() string {
	s, _ := e.clipboard.Get()

	return s
}

// SetClipboard modifies the contents of the clipboard. It assumes the text was
// not captured from an auto-expanded cursor. It exists for Sublime Text API
// compatibility.
func (e *Editor) SetClipboard(s string) {
	e.clipboard.Set(s, false)
}

func (e *Editor) handleLog(s string) {
	c := e.Console()
	f := fmt.Sprintf("%08d %d %s", c.Size(), len(s), s)
	edit := c.BeginEdit()
	c.Insert(edit, c.Size(), f)
	c.EndEdit(edit)
}

func (e *Editor) AddPackagesPath(p string) {
	e.pkgsPaths = append(e.pkgsPaths, p)
	OnPackagesPathAdd.call(p)
}

func (e *Editor) RemovePackagesPath(p string) {
	e.pkgsPaths = util.Remove(e.pkgsPaths, p)
	OnPackagesPathRemove.call(p)
}

func (e *Editor) DefaultPath() string {
	return e.defaultPath
}

func (e *Editor) UserPath() string {
	return e.userPath
}

func (e *Editor) SetDefaultPath(p string) {
	e.defaultPath = p
	OnDefaultPathAdd.call(p)
}

func (e *Editor) SetUserPath(p string) {
	e.userPath = p
	OnUserPathAdd.call(p)
}

func (e *Editor) AddColorScheme(path string, cs ColorScheme) {
	e.colorSchemes[path] = cs
}

func (e *Editor) GetColorScheme(path string) ColorScheme {
	if path == "" {
		return defaultScheme()
	}

	scheme := e.colorSchemes[path]
	if scheme == nil {
		log.Error("No color scheme %s in editor falling back to default color scheme", path)
		return defaultScheme()
	}
	return scheme
}

// TODO: should generate sth like sublime text color schemes menu
func (e *Editor) ColorSchemes() {}

func (e *Editor) AddSyntax(path string, s Syntax) {
	e.syntaxes[path] = s
	for _, t := range s.FileTypes() {
		e.filetypes[t] = path
	}
}

func (e *Editor) GetSyntax(path string) Syntax {
	return e.syntaxes[path]
}

func (e *Editor) fileTypeSyntax(ext string) string {
	return e.filetypes[ext]
}

// TODO: should generate sth like sublime text syntaxes menu
// the name in the menu should come from defined name inside syntax file or
// the syntax file name
func (e *Editor) Syntaxes() {}
