package manga_search

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	path := "../data/manga.txt"

	db, err := Load(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	r := db.Search("死")
	for _, m := range r {
		fmt.Printf("%+v\n", m)
	}

}
