package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCompress(t *testing.T) {
	txt := `// AllocsPerRun returns the average number of allocations during calls to f.
// Although the return value has type float64, it will always be an integral value.
//
// To compute the number of allocations, the function will first be run once as
// a warm-up. The average number of allocations over the specified number of
// runs will then be measured and returned.
//
// AllocsPerRun sets GOMAXPROCS to 1 during its measurement and will restore
// it before returning.
`
	txtInput := strings.NewReader(txt)
	img := bytes.NewBuffer(nil)
	err := compress(txtInput, img)
	if err != nil {
		t.Fatal("compress err", err)
	}
	txtOutput := strings.Builder{}
	err = uncompress(img, &txtOutput)
	if err != nil {
		t.Fatal("uncompress err", err)
	}
	if txtOutput.String() != txt {
		t.Fatal("uncompressed text does not match original text")
	}
}
