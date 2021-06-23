package clipboard

import (
	"fmt"
	"os"
	"testing"
)

func TestParseStoreFile(t *testing.T) {
	f, err := os.Open("testfile1.txt")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	marks, title := ParseStoreFile(f)
	if title != "202102" {
		t.Error(title)
	}
	for _, m := range marks {
		fmt.Println(m.String())
	}
	if len(marks) != 2 {
		t.Error(len(marks))
	}
}

func TestMerge(t *testing.T) {
	f, err := os.Open("testfile1.txt")
	if err != nil {
		t.Error(err)
	}
	f2, err := os.Open("testfile2.txt")
	if err != nil {
		t.Error(err)
	}
	defer f2.Close()
	marks, _ := ParseStoreFile(f)

	marks2, _ := ParseStoreFile(f2)

	mergeMark := Merge(marks,marks2)
	for _, m := range mergeMark {
		fmt.Println(m.String())
	}
	if len(mergeMark) != 3 {
		t.Error(len(mergeMark))
	}
}
