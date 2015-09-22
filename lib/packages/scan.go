package packages

import (
	"io/ioutil"
	"path"

	"github.com/limetext/lime-backend/lib/log"
)

// A helper struct to watch all scaned directories for new packages
type scanDir struct {
	path string
}

func (p *scanDir) FileCreated(name string) {
	record(name)
}

// watches scaned directory
func watchDir(dir string) {
	sd := &scanDir{dir}
	if err := watcher.Watch(sd.path, sd); err != nil {
		log.Warn("Couldn't watch %s: %s", sd.path, err)
	}
}

func Scan(dir string) {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Warn("Couldn't read path %s: %s", dir, err)
	}

	watchDir(dir)

	for _, fi := range fis {
		record(path.Join(dir, fi.Name()))
	}
}
