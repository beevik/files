// Package files contains functions and types that help with file and
// directory processing.
package files

import (
    "errors"
    "os"
    "path"
    "regexp"
)

// A DirMode represents flags that define the behavior of the DirReader.
type DirMode uint32

const (
    D_RECURSE DirMode = 1 << iota // Recurse all subdirectories
)

// Errors returned by the files packages
var (
    ErrNotDirectory error
)

// A DirReader iterates through the files and subdirectories contained
// within a directory.
type DirReader struct {
    Filter        Filter
    mode          DirMode
    dirs          []FileInfo
    files         []FileInfo
    dirsTraversed int
}

// A FileInfo contains a file's full path.  It also embeds an os.FileInfo.
type FileInfo struct {
    Path string
    os.FileInfo
}

func init() {
    ErrNotDirectory = errors.New("files: file is not a directory")
}

// A Filter interface is used to define rules to include or exclude
// a file in the results of a Next() iteration.
type Filter interface {
    Eval(f *FileInfo) bool
}

type fileFilter struct{}

func (ff fileFilter) Eval(f *FileInfo) bool {
    return !f.IsDir()
}

// FileFilter creates a filter that accepts only files (not directories).
func FileFilter() Filter {
    return fileFilter{}
}

type dirFilter struct{}

func (df dirFilter) Eval(f *FileInfo) bool {
    return f.IsDir()
}

// DirFilter creates a filter that accepts only directories (not files).
func DirFilter() Filter {
    return dirFilter{}
}

type regexpFilter struct {
    pattern *regexp.Regexp
}

func (rf regexpFilter) Eval(f *FileInfo) bool {
    return rf.pattern.MatchString(f.Path)
}

// RegexpFilter creates a filter that returns true when the regular
// expression matches the file's full path.
func RegexpFilter(pattern string) Filter {
    p := regexp.MustCompile(pattern)
    return regexpFilter{p}
}

type multiFilter struct {
    filters []Filter
}

func (mf multiFilter) Eval(f *FileInfo) bool {
    for _, fx := range mf.filters {
        if !fx.Eval(f) {
            return false
        }
    }
    return true
}

// MultiFilter creates a filter composed of several other filters.
func MultiFilter(filters ...Filter) Filter {
    return &multiFilter{filters}
}

// NewDirReader creates a DirReader rooted at the specified
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
        return nil, ErrNotDirectory
    }

    r := &DirReader{
        dirs: make([]FileInfo, 1, 8),
        mode: mode,
    }
    r.dirs[0] = FileInfo{dir, info}
    return r, nil
}

// Next iterates to the next available file in the directory and
// returns its FileInfo.  Returns a nil FileInfo when complete.
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

        // Examine the next file
        var info *FileInfo
        info, r.files = &r.files[0], r.files[1:]
        if (r.mode&D_RECURSE) == D_RECURSE && info.IsDir() {
            r.dirs = append(r.dirs, *info)
        }
        if r.Filter == nil || r.Filter.Eval(info) {
            return info, nil
        }
    }
    return nil, nil
}

// getMoreFiles is a helper function that retrieves more files
// from a directory.
func (r *DirReader) getMoreFiles() error {
    var info *FileInfo
    info, r.dirs = &r.dirs[0], r.dirs[1:]
    if r.dirsTraversed++; r.dirsTraversed%64 == 0 {
        newdirs := make([]FileInfo, len(r.dirs))
        copy(newdirs, r.dirs)
        r.dirs = newdirs
    }

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
        r.files = append(r.files, FileInfo{path.Join(info.Path, i.Name()), i})
    }
    return nil
}
