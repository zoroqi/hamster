package main

import (
	"testing"
)

func TestRename(t *testing.T) {
	mapping := defaultMapping()
	testcase := map[string]string{
		"ab":         "ab",
		"a b":        "a-b",
		"ab ":        "ab",
		"ab;;":       "ab",
		"ab.txt":     "ab.txt",
		"a b.txt":    "a-b.txt",
		"a b .txt":   "a-b.txt",
		"a:::b.txt":  "a-b.txt",
		"a;;b;;.txt": "a-b.txt",
		";;;ab.txt":  "ab.txt",
		";;;ab":      "ab",
		"ab\\ab":     "ab-ab",
		".abc":       ".abc",
		". ab":       ". ab",
	}
	for k, v := range testcase {
		if rename(k, mapping) != v {
			t.Error(k, "->", v, " but ", rename(k, mapping))
		}
	}
}
