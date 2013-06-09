package files

import (
    "testing"
)

func TestNewDirReader(t *testing.T) {
    r, err := NewDirReader("C:/Temp", D_RECURSE)
    if err != nil {
        panic(err)
    }
    r.Filter = MultiFilter(FileFilter(), RegexpFilter(`\.xml$`))
    for {
        info, err := r.Next()
        if err != nil {
            panic(err)
        }
        if info == nil {
            break
        }
    }
}
