// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package watch

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func newWatcher(t *testing.T) *Watcher {
	watcher, err := NewWatcher()
	if err != nil {
		t.Fatalf("Couldn't create watcher: %s", err)
	}
	return watcher
}

func watch(t *testing.T, watcher *Watcher, name string, cb interface{}) {
	if err := watcher.Watch(name, cb); err != nil {
		t.Fatalf("Couldn' Watch %s : %s", name, err)
	}
}

func unwatch(t *testing.T, watcher *Watcher, name string, cb interface{}) {
	if err := watcher.UnWatch(name, cb); err != nil {
		t.Fatalf("Couldn' UnWatch %s : %s", name, err)
	}
}

func testWatched(t *testing.T, watched map[string][]interface{}, expWatched []string) {
	if len(watched) != len(expWatched) {
		t.Errorf("Expected watched %v keys equal to %v", watched, expWatched)
	}
	for _, p := range expWatched {
		absp, err := filepath.Abs(p)
		if err != nil {
			t.Errorf("Failed to Abs(%s): %s", p, err)
		}
		if _, exist := watched[absp]; !exist {
			t.Errorf("Expected %s exist in watched", absp)
		}
	}
}

func testWatchers(t *testing.T, watchers []string, expWatchers []string) {
	if len(watchers) != len(expWatchers) {
		t.Errorf("Expected watchers %v keys equal to %v", watchers, expWatchers)
	}
	for i, p := range expWatchers {
		absp, err := filepath.Abs(p)
		if err != nil {
			t.Errorf("Failed to Abs(%s): %s", p, err)
		}
		if watchers[i] != absp {
			t.Errorf("Expected watchers %s to be %s", watchers[i], absp)
		}
	}
}

func TestNewWatcher(t *testing.T) {
	watcher := newWatcher(t)
	defer watcher.Close()
	if len(watcher.dirs) != 0 {
		t.Errorf("Expected len(dirs) of new watcher %d, but got %d", 0, len(watcher.dirs))
	}
	if len(watcher.watchers) != 0 {
		t.Errorf("Expected len(watchers) of new watcher %d, but got %d", 0, len(watcher.watchers))
	}
}

type dummy struct {
	name    string
	c       chan bool
	lock    sync.Mutex
	created bool
	changed bool
	renamed bool
	removed bool
}

func newDummy(name string) *dummy {
	return &dummy{name: name, c: make(chan bool, 5)}
}

func (d *dummy) reset() {
	d.created = false
	d.changed = false
	d.renamed = false
	d.removed = false
}

func (d *dummy) done(name string, got *bool) {
	// fmt.Println("Dummy: ", got, name == d.name, name, d.name)
	if name != d.name {
		return
	}
	d.lock.Lock()
	defer func() { d.c <- true }() // make sure Unlock() is called first
	defer d.lock.Unlock()          // in order to avoid deadlocks
	*got = true
}

func (d *dummy) FileChanged(name string) {
	d.done(name, &d.changed)
}

func (d *dummy) FileCreated(name string) {
	d.done(name, &d.created)
}

func (d *dummy) FileRemoved(name string) {
	d.done(name, &d.removed)
}

func (d *dummy) FileRenamed(name string) {
	d.done(name, &d.renamed)
}

func (d *dummy) Wait() {
	<-d.c
	for {
		select {
		case <-d.c:
			continue
		case <-time.After(10 * time.Millisecond):
			return
		}
	}
}

func TestWatch(t *testing.T) {
	tests := []struct {
		paths       []string
		expWatched  []string
		expWatchers []string
	}{
		{
			[]string{"testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata/dummy.txt", "testdata/test.txt"},
		},
		{
			[]string{"testdata", "testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata", "testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata"},
		},
		{
			[]string{"testdata/dummy.txt", "testdata/test.txt", "testdata"},
			[]string{"testdata", "testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata"},
		},
	}
	for _, test := range tests {
		watcher := newWatcher(t)
		for _, name := range test.paths {
			watch(t, watcher, name, newDummy(name))
		}
		testWatched(t, watcher.watched, test.expWatched)
		testWatchers(t, watcher.watchers, test.expWatchers)
		defer watcher.Close()
	}
}

func Testwatch(t *testing.T) {
	watcher := newWatcher(t)
	defer watcher.Close()
	if err := watcher.watch("testdata/dummy.txt", false); err != nil {
		t.Fatalf("Couldn't watch %s", "testdata/dummy.txt")
	}
	if err := watcher.watch("testdata/test.txt", false); err != nil {
		t.Fatalf("Couldn't watch %s", "testdata/test.txt")
	}
	testWatched(t, watcher.watched, []string{"testdata/dummy.txt", "testdata/test.txt"})
	testWatchers(t, watcher.watchers, []string{"testdata/dummy.txt", "testdata/test.txt"})
}

func TestAdd(t *testing.T) {
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy("test")
	watcher.add("test", d)
	if cb := watcher.watched["test"][0]; cb != d {
		t.Errorf("Expected watcher['test'][0] callback equal to %v, but got %v", d, cb)
	}
}

func TestFlushDir(t *testing.T) {
	name := "testdata/dummy.txt"
	dir, _ := filepath.Abs("testdata")
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy(name)
	watch(t, watcher, name, d)
	testWatchers(t, watcher.dirs, []string{})
	testWatchers(t, watcher.watchers, []string{name})
	watcher.flushDir(dir)
	testWatchers(t, watcher.dirs, []string{dir})
	testWatchers(t, watcher.watchers, []string{})
}

func TestUnWatch(t *testing.T) {
	name, _ := filepath.Abs("testdata/dummy.txt")
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy(name)
	watch(t, watcher, name, d)
	unwatch(t, watcher, name, d)
	if len(watcher.watched) != 0 {
		t.Errorf("Expected watcheds be empty, but got %v", watcher.watched)
	}
}

func TestUnWatchAll(t *testing.T) {
	name, _ := filepath.Abs("testdata/dummy.txt")
	watcher := newWatcher(t)
	defer watcher.Close()
	d1 := new(dummy)
	d2 := new(dummy)
	watch(t, watcher, name, d1)
	watch(t, watcher, name, d2)
	if l := len(watcher.watched[name]); l != 2 {
		t.Errorf("Expected len of watched['%s'] be %d, but got %d", name, 2, l)
	}
	unwatch(t, watcher, name, nil)
	if _, exist := watcher.watched[name]; exist {
		t.Errorf("Expected all %s watched be removed", name)
	}
	testWatchers(t, watcher.watchers, []string{})
}

func TestUnWatchDirectory(t *testing.T) {
	name := "testdata/dummy.txt"
	absname, _ := filepath.Abs(name)
	dir, _ := filepath.Abs("testdata")
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy(absname)
	watch(t, watcher, name, d)
	watch(t, watcher, dir, nil)
	testWatchers(t, watcher.watchers, []string{dir})
	unwatch(t, watcher, dir, nil)
	testWatchers(t, watcher.watchers, []string{name})
}

func TestUnWatchOneOfSubscribers(t *testing.T) {
	name, _ := filepath.Abs("testdata/dummy.txt")
	watcher := newWatcher(t)
	defer watcher.Close()
	d1 := new(dummy)
	d2 := new(dummy)
	watch(t, watcher, name, d1)
	watch(t, watcher, name, d2)
	if len(watcher.watched[name]) != 2 {
		t.Fatalf("Expected watched[%s] length be %d, but got %d", name, 2, len(watcher.watched[name]))
	}
	unwatch(t, watcher, name, d1)
	testWatchers(t, watcher.watchers, []string{name})
	if len(watcher.watched[name]) != 1 {
		t.Errorf("Expected watched[%s] length be %d, but got %d", name, 1, len(watcher.watched[name]))
	}
}

func TestunWatch(t *testing.T) {
	name := "testdata/dummy.txt"
	watcher := newWatcher(t)
	defer watcher.Close()
	d1 := new(dummy)
	d2 := new(dummy)
	watch(t, watcher, name, d1)
	watch(t, watcher, name, d2)
	if err := watcher.unWatch(name); err != nil {
		t.Fatalf("Couldn't unWatch %s: %s", name, err)
	}
	if _, exist := watcher.watched[name]; exist {
		t.Errorf("Expected all %s watched be removed", name)
	}
	// if !reflect.DeepEqual(watcher.watchers, []string{}) {
	// 	t.Errorf("Expected watchers be empty but got %v", watcher.watchers)
	// }
	testWatchers(t, watcher.watchers, []string{})
}

func TestRemoveWatch(t *testing.T) {
	name, _ := filepath.Abs("testdata/dummy.txt")
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy(name)
	watch(t, watcher, name, d)
	watcher.removeWatch(name)
	// if !reflect.DeepEqual(watcher.watchers, []string{}) {
	// 	t.Errorf("Expected watchers be empty but got %v", watcher.watchers)
	// }
	testWatchers(t, watcher.watchers, []string{})
}

func TestRemoveDir(t *testing.T) {
	name, _ := filepath.Abs("testdata/dummy.txt")
	dir, _ := filepath.Abs("testdata")
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy(name)
	watch(t, watcher, dir, d)
	watch(t, watcher, name, d)
	testWatchers(t, watcher.watchers, []string{dir})
	testWatchers(t, watcher.dirs, []string{dir})
	watcher.removeDir(dir)
	testWatchers(t, watcher.dirs, []string{})
	testWatchers(t, watcher.watchers, []string{dir, name})
}

func TestObserve(t *testing.T) {
	name := "testdata/test.txt"
	absname, _ := filepath.Abs(name)
	watcher := newWatcher(t)
	defer ioutil.WriteFile(name, []byte(""), 0644)
	defer watcher.Close()
	d := newDummy(absname)
	watch(t, watcher, name, d)

	if err := ioutil.WriteFile(name, []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile error: %s", err)
	}

	d.Wait()
	if !d.changed {
		t.Errorf("Expected dummy Text %s, but got %#v", "Changed", d)
	}
}

func TestCreateEvent(t *testing.T) {
	name := "testdata/new.txt"
	absname, _ := filepath.Abs(name)
	os.Remove(name)
	defer os.Remove(name)
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy(absname)
	watch(t, watcher, name, d)

	testWatchers(t, watcher.watchers, []string{"testdata"})

	if f, err := os.Create(name); err != nil {
		t.Fatalf("File creation error: %s", err)
	} else {
		f.Close()
	}
	d.Wait()
	if !d.created {
		t.Errorf("Expected dummy Text %s, but got %#v", "Created", d)
	}
}

func TestDeleteEvent(t *testing.T) {
	if os.ExpandEnv("$TRAVIS") != "" {
		// This test just times out on travis (ie the callback is never called).
		// See https://github.com/limetext/lime/issues/438
		t.Skip("Skipping test as it doesn't work with travis")
		return
	}
	name := "testdata/dummy.txt"
	absname, _ := filepath.Abs(name)
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy(absname)
	watch(t, watcher, name, d)

	if err := os.Remove(name); err != nil {
		t.Fatalf("Couldn't remove file %s: %s", name, err)
	}
	d.Wait()
	if !d.removed {
		t.Errorf("Expected dummy Text %s, but got %#v", "Removed", d)
	}
	if f, err := os.Create(name); err != nil {
		t.Errorf("Couldn't create file: %s", err)
	} else {
		f.Close()
	}
	d.Wait()
	if !d.created {
		t.Errorf("Expected dummy Text %s, but got %#v", "Created", d)
	}
}

func TestRenameEvent(t *testing.T) {
	name := "testdata/test.txt"
	absname, _ := filepath.Abs(name)
	defer os.Rename("testdata/rename.txt", name)
	watcher := newWatcher(t)
	defer watcher.Close()
	d := newDummy(absname)
	watch(t, watcher, name, d)

	os.Rename(name, "testdata/rename.txt")
	d.Wait()
	if !d.renamed {
		t.Errorf("Expected dummy Text %s, but got %#v", "Renamed", d)
	}
}
