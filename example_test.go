package files_test

import (
	"fmt"

	"github.com/beevik/files"
)

// Create a recursive directory reader on the directory C:/Temp and
// output all files ending in the extension .xml.
func ExampleNewDirReader() {
	r, err := files.NewDirReader("C:/Temp", files.Recurse)
	if err != nil {
		fmt.Println("ERROR", err)
		return
	}
	r.Filter = files.MultiFilter(files.FileFilter(), files.RegexpFilter(`\.xml$`))
	for {
		info, err := r.Next()
		if err != nil {
			fmt.Println("ERROR", err)
			return
		}
		if info == nil {
			break
		}
		fmt.Println(info.Path)
	}
}
