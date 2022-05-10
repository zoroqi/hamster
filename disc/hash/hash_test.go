package hash

import (
	"testing"
)

func TestFileFastMd5sum(t *testing.T) {
	sum, err := FileFastMd5sum("../develop_log.md", 1024)
	if err != nil {
		t.Fatal(err)
	}
	if sum != "3f6142726b075242a26f78d90bf64832" {
		t.Fatal(sum, "!=", "3f6142726b075242a26f78d90bf64832")
	}
}
