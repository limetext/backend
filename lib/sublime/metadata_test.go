package sublime

import (
	"io/ioutil"
	"testing"

	"github.com/quarnster/completion/util"
)

func TestLoadMetadata(t *testing.T) {
	var (
		in  = "testdata/Comments.tmPreferences"
		exp = "testdata/Comments.tmPreferences.res"
	)

	md, err := LoadMetadata(in)
	if err != nil {
		t.Fatalf("Error on loading %s: %s", in, err)
	}
	data, err := ioutil.ReadFile(exp)
	if err != nil {
		t.Fatalf("Error reading expected file %s: %s", exp, err)
	}
	if diff := util.Diff(string(data), md.String()); diff != "" {
		t.Error(diff)
	}
}
