package files

import (
	"os"
	"testing"
)

func TestNewDirReader(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	r, err := NewDirReader(wd, Recurse)
	if err != nil {
		t.Error(err)
	}
	r.Filter = MultiFilter(FileFilter(), RegexpFilter(`\.go$`))
	for {
		info, err := r.Next()
		if err != nil {
			t.Error(err)
		}
		if info == nil {
			break
		}
	}
}
