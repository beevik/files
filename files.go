package files

import (
    "errors"
    "os"
    "regexp"
)

// A DirMode represents flags that define the behavior of the DirReader.
type DirMode uint32

const (
    D_RECURSE DirMode = 1 << iota   // Recurse all directories
)

// A DirReader iterates through the files contained within a directory.
type DirReader struct {
    Filter Filter
    mode   DirMode
    dirs   []FileInfo
    files  []FileInfo
}

// A FileInfo contains a file's full path.  It also embeds a os.FileInfo struct.
type FileInfo struct {
    Path string
    os.FileInfo
}

// A Filter defines rules that allow inclusion or exclusion of a file
// in the results of a DirReader iteration.
type Filter interface {
    Test(f *FileInfo) bool
}

type fileFilter struct{}

func (ff fileFilter) Test(f *FileInfo) bool {
    return !f.IsDir()
}

// Create a filter that returns true when passed a file (not a directory).
func FileFilter() Filter {
    return fileFilter{}
}

type dirFilter struct{}

func (df dirFilter) Test(f *FileInfo) bool {
    return f.IsDir()
}

// Createa a filter that returns true when passed a directory.
func DirFilter() Filter {
    return dirFilter{}
}

type regexpFilter struct {
    pattern *regexp.Regexp
}

func (rf regexpFilter) Test(f *FileInfo) bool {
    return rf.pattern.MatchString(f.Path)
}

// Create a filter that returns true when the regular expression matches
// the file's full path.
func RegexpFilter(pattern string) Filter {
    p := regexp.MustCompile(pattern)
    return regexpFilter{p}
}

type multiFilter struct {
    filters []Filter
}

func (mf multiFilter) Test(f *FileInfo) bool {
    for _, fx := range mf.filters {
        if !fx.Test(f) {
            return false
        }
    }
    return true
}

// Create a filter composed of several other filters.
func MultiFilter(filters ...Filter) Filter {
    return &multiFilter{filters}
}

// NewDirReader creates a new directory reader rooted at the specified
// directory.
func NewDirReader(dir string, mode DirMode) (*DirReader, error) {
    f, err := os.OpenFile(dir, os.O_RDONLY, os.ModePerm)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    info, err := f.Stat()
    if err != nil {
        return nil, err
    }
    if !info.IsDir() {
        return nil, errors.New("file is not a directory")
    }

    r := new(DirReader)
    r.dirs = make([]FileInfo, 1, 8)
    r.dirs[0] = FileInfo{dir, info}
    r.mode = mode
    return r, nil
}

// Next returns the next file in the iteration of the directory
func (r *DirReader) Next() (*FileInfo, error) {
    for {
        // Retrieve more files if available
        for len(r.files) == 0 {
            if len(r.dirs) == 0 {
                return nil, nil
            }
            if err := r.getMoreFiles(); err != nil {
                return nil, err
            }
        }

        // Test the next file
        var info FileInfo
        info, r.files = r.files[0], r.files[1:]
        if (r.mode&D_RECURSE) == D_RECURSE && info.IsDir() {
            r.dirs = append(r.dirs, info)
        }
        if r.Filter == nil || r.Filter.Test(&info) {
            return &info, nil
        }
    }
    return nil, nil
}

// getMoreFiles is a hlper function that retrieves more files
// from a directory
func (r *DirReader) getMoreFiles() error {
    var info FileInfo
    info, r.dirs = r.dirs[0], r.dirs[1:]

    f, err := os.OpenFile(info.Path, os.O_RDONLY, os.ModePerm)
    if err != nil {
        return err
    }
    defer f.Close()

    files, err := f.Readdir(0)
    if err != nil {
        return err
    }
    for _, i := range files {
        r.files = append(r.files, FileInfo{info.Path + "/" + i.Name(), i})
    }
    return nil
}
