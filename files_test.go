package files

import (
    "testing"
)

func TestNewDirReader(t *testing.T) {
    r, err := NewDirReader("C:/Temp", D_RECURSE)
    //r.Filter = MultiFilter(DirFilter(), RegexpFilter(`_`))
    r.Filter = DirFilter()

    if err != nil {
        t.Error(err)
        return
    }
    for {
        i, err := r.Next()
        if err != nil {
            t.Error(err)
        }
        if i == nil {
            break
        }
        t.Logf("%v %v", string(typechar(i)), i.Path)
    }
}

func typechar(i *FileInfo) byte {
	if i.IsDir() {
		return 'D'
	}
	return 'F'
}