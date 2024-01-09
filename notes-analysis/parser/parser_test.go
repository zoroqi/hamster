package parser

import (
	"fmt"
	"path/filepath"
	"testing"
)

func Test_GoldMark(t *testing.T) {
	//path := "_data/first.md"
	//f, err := os.Open(path)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//ParseMarkdown(f, path)
	fmt.Println(filepath.Join("nofile", "abc.md"))
	//	fmt.Printf("%q", `\ ab"
	//c`)
	//
	//fmt.Println(root.Lines().Len())
	//
	//metadata := meta.GetItems(context)
	//var metadfs func(any, int) bool
	//metadfs = func(a any, dep int) bool {
	//	//fmt.Println(reflect.TypeOf(a))
	//	switch vt := a.(type) {
	//	case yaml.MapSlice:
	//		for _, v := range vt {
	//			metadfs(v, dep+1)
	//		}
	//		return false
	//	case yaml.MapItem:
	//		fmt.Printf("%s%s:\n", strings.Repeat(" ", dep), vt.Key)
	//		metadfs((vt.Value), dep+1)
	//		return false
	//	case []any:
	//		for _, v := range vt {
	//			metadfs(v, dep+1)
	//		}
	//		return false
	//	default:
	//		fmt.Printf("%s%v\n", strings.Repeat(" ", dep), vt)
	//		return true
	//	}
	//}
	//
	//metadfs(metadata, 0)
}
