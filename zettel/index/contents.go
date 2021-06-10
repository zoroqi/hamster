package index

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func ZettelIndex() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index",
		Long:  "build zettel index [给每条笔记分配一个固定地址](https://zettelkasten.de/introduction/zh/#%E7%BB%99%E6%AF%8F%E6%9D%A1%E7%AC%94%E8%AE%B0%E5%88%86%E9%85%8D%E4%B8%80%E4%B8%AA%E5%9B%BA%E5%AE%9A%E5%9C%B0%E5%9D%80)",
		Short: "build zettel index",

		Run: func(cmd *cobra.Command, args []string) {
			dir, _ := cmd.Flags().GetString("directory")
			output, _ := cmd.Flags().GetString("output")
			pre, _ := cmd.Flags().GetString("pre")
			buildIndex(dir, output, pre)
		},
	}
	cmd.Flags().StringP("directory", "d", "", "scan directory")
	cmd.Flags().StringP("output", "o", "", "file output path")
	cmd.Flags().StringP("pre", "p", "", "scan file name prefix")
	return cmd
}

type zettelIndex struct {
	name  []string
	title string
	path  string
	child map[string]zettelIndex
}

func buildZI(name string, pre string, path string) zettelIndex {
	if len(pre) > 0 {
		name = name[len(pre):]
	}
	// 移出.md
	name = name[:len(name)-3]
	name = strings.ToLower(name)
	numRegex := regexp.MustCompile("\\d+")
	letterRegex := regexp.MustCompile("[a-z]+")
	nums := numRegex.FindAllString(name, -1)
	letters := letterRegex.FindAllString(name, -1)

	z := zettelIndex{path: path, name: make([]string, 0, len(nums)+len(letters)), child: make(map[string]zettelIndex, 0)}

	l := len(nums)
	if l < len(letters) {
		l = len(letters)
	}
	for i := 0; i < l; i++ {
		if i < len(nums) {
			z.name = append(z.name, nums[i])
		}
		if i < len(letters) {
			z.name = append(z.name, letters[i])
		}
	}

	txt, err := ioutil.ReadFile(path)
	if err == nil {
		p := regexp.MustCompile("^#\\s+(.*?)\n")
		title := p.FindAllSubmatch(txt, -1)
		if len(title) > 0 {
			z.title = string(title[0][1])
		}
	}

	return z
}

func (z zettelIndex) Name() string {
	return strings.Join(z.name, "")
}
func (z zettelIndex) FileName() string {
	return strings.Join(z.name, "") + ".md"
}

func toString(z zettelIndex, root string) string {
	sb := strings.Builder{}
	var dfs func(z zettelIndex, root string, space string)
	dfs = func(z zettelIndex, root string, space string) {
		var p string
		if strings.HasPrefix(z.path, root) {
			p = strings.Replace(z.path, root, ".", 1)
		} else {
			p = fmt.Sprintf("%s%s", root, z.path)
		}
		sb.WriteString(fmt.Sprintf("%s* [%s|%s](%s)\n", space, z.Name(), z.title, p))
		arr := make([]zettelIndex, len(z.child))
		i := 0
		for _, c := range z.child {
			arr[i] = c
			i++
		}
		sortZettel(arr)
		for _, c := range arr {
			dfs(c, root, space+"    ")
		}
	}
	dfs(z, root, "")
	return sb.String()
}

// 1 < 2, 1a < 1b, 越短越大
func sortZettel(arr []zettelIndex) {
	sort.Slice(arr, func(i, j int) bool {
		in := arr[i].name
		jn := arr[j].name

		l := len(in)
		if l > len(jn) {
			l = len(jn)
		}

		for i := 0; i < l; i++ {
			if in[i] != jn[i] {
				return in[i] < jn[i]
			}
		}
		return len(in) < len(jn)
	})
}

func buildIndex(dir, output, pre string) {
	if dir == "" {
		fmt.Println("must dir")
		return
	}
	if output == "" {
		output = "./"
	}
	files := make([]zettelIndex, 0)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return err
		}
		if strings.HasSuffix(info.Name(), ".md") && strings.HasPrefix(info.Name(), pre) {
			z := buildZI(info.Name(), pre, path)
			files = append(files, z)
		}
		return err
	})
	if err != nil {
		fmt.Errorf("walk err %s", err)
		return
	}
	sortZettel(files)
	files = merge(files)
	for _, f := range files {
		ioutil.WriteFile(
			filepath.Join(output, fmt.Sprintf("index%s.md", f.name[0])),
			[]byte(toString(f, dir)),
			0666)
	}

}

func merge(files []zettelIndex) []zettelIndex {

	index := make(map[string]zettelIndex)

	for _, f := range files {
		in := index
		for _, ns := range f.name {
			if len(in[ns].path) == 0 {
				in[ns] = f
				break
			} else {
				in = in[ns].child
			}
		}
	}

	r := make([]zettelIndex, 0, len(index))
	for _, v := range index {
		r = append(r, v)
	}
	return r
}
