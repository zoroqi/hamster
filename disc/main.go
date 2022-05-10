// +build ignore

package main

import "fmt"
import "github.com/zoroqi/rubbish/disc"

func main() {
	meta := disc.FastScan(".")
	meta2, _ := disc.SlowScan(meta)
	for _, m := range meta2 {
		fmt.Println(m)
	}
}
