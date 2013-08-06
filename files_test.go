package files

import (
    "io/ioutil"
    "os"
    "testing"
)

func TestNewDirReader(t *testing.T) {
    curdir, _ := os.Getwd()
    r, err := NewDirReader(curdir, D_RECURSE)
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

func TestGetAllFilePaths(t *testing.T) {
    tempDir := os.TempDir()
    r1, _ := NewDirReader(tempDir, D_RECURSE)
    numFiles1 := len(r1.GetAllFilePaths())
    ioutil.TempFile(tempDir, "foo")
    r2, _ := NewDirReader(tempDir, D_RECURSE)
    numFiles2 := len(r2.GetAllFilePaths())
    if numFiles2 != numFiles1+1 {
        t.Errorf("First try: %d, Second try: %d", numFiles1, numFiles2)    
    }
}
