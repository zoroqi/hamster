package file

import (
	"archive/zip"
	"fmt"
)

func Unziplv(path2 string) {
	r, err := zip.OpenReader(path2)
	if err != nil {
		return
	}

	for _, f := range r.File {
		fmt.Printf("%+v\n", f.FileHeader)
	}
}