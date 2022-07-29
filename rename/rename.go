package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type info struct {
	path     string
	fileInfo os.FileInfo
}

func main() {
	args := os.Args
	var root string
	if len(args) == 1 {
		root = "."
	} else {
		root = args[1]
	}
	renameDir(root)
	renameFile(root)
}

func renameDir(root string) {
	queue := make([]struct {
		p string
	}, 0)
	queue = append(queue, struct{ p string }{p: root})
	bfs := func(path string) error {
		fileInfo, err := os.Lstat(path)
		if err != nil {
			return err
		}
		newpath := path
		if fileInfo.IsDir() && !strings.HasPrefix(fileInfo.Name(), ".") {
			i := info{path: path, fileInfo: fileInfo}
			newpath, err = replaceName(i)
			if err != nil {
				return err
			}

		}
		fileInfo, err = os.Lstat(newpath)
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			files, err := readDirNames(newpath)
			if err != nil {
				return err
			}
			for _, f := range files {
				nfp := filepath.Join(newpath, f)
				fi, err := os.Lstat(nfp)
				if err != nil {
					return err
				}
				if fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") {
					queue = append(queue, struct {
						p string
					}{p: nfp})
				}
			}
		}
		return nil
	}

	for len(queue) != 0 {
		if err := bfs(queue[0].p); err != nil {
			fmt.Println(err)
			return
		}
		queue = queue[1:]
	}
}

func readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

func replaceNames(infos []info) error {
	for _, info := range infos {
		_, err := replaceName(info)
		if err != nil {
			return err
		}
	}
	return nil
}

func replaceName(i info) (newPath string, err error) {
	mapping := loadMapping()
	if mapping == nil {
		return i.path, nil
	}
	newName := rename(i.fileInfo.Name(), mapping)
	if newName != i.fileInfo.Name() {
		fmt.Println(i.fileInfo.Name() + " -> " + newName)
		path := i.path
		dir, _ := filepath.Split(path)
		newPath := filepath.Join(dir, newName)
		err := os.Rename(i.path, newPath)
		return newPath, err
	}
	return i.path, err
}

func renameFile(root string) {
	var infos []info
	err := filepath.Walk(root, func(path string, fileInfo os.FileInfo, err error) error {
		if err == nil && !fileInfo.IsDir() && !strings.HasPrefix(fileInfo.Name(), ".") {
			infos = append(infos, info{path: path, fileInfo: fileInfo})
		}
		if root != path && fileInfo.IsDir() && strings.HasPrefix(fileInfo.Name(), ".") {
			return filepath.SkipDir
		}
		return err
	})
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	replaceNames(infos)
}

func loadMapping() map[rune]bool {
	return defaultMapping
}

var defaultMapping = make(map[rune]bool)

func init() {
	defaultMapping[' '] = true
	defaultMapping['　'] = true

	defaultMapping['：'] = true
	defaultMapping[':'] = true

	defaultMapping['；'] = true
	defaultMapping[';'] = true

	defaultMapping['，'] = true
	defaultMapping['。'] = true
	defaultMapping['、'] = true

	defaultMapping['《'] = true
	defaultMapping['》'] = true

	defaultMapping['—'] = true

	defaultMapping['（'] = true
	defaultMapping['）'] = true

	defaultMapping['+'] = true

	defaultMapping['【'] = true
	defaultMapping['】'] = true
	defaultMapping['|'] = true
	defaultMapping['！'] = true
	defaultMapping['!'] = true
	defaultMapping['／'] = true
	defaultMapping['\\'] = true
	defaultMapping['〜'] = true
	defaultMapping['~'] = true
	defaultMapping['\''] = true
	defaultMapping['"'] = true

	defaultMapping['’'] = true
	defaultMapping['‘'] = true
	defaultMapping['”'] = true
	defaultMapping['“'] = true
	defaultMapping['…'] = true
	defaultMapping['*'] = true
	defaultMapping[','] = true
	defaultMapping[')'] = true
	defaultMapping['('] = true
	defaultMapping['['] = true
	defaultMapping[']'] = true
	defaultMapping['&'] = true
	defaultMapping['#'] = true
	defaultMapping['^'] = true
	defaultMapping['`'] = true
	defaultMapping['@'] = true
	defaultMapping['・'] = true
	defaultMapping['{'] = true
	defaultMapping['}'] = true
	defaultMapping['「'] = true
	defaultMapping['」'] = true
}

func rename(name string, mapping map[rune]bool) string {
	const midlineS = "-"
	const midlineR = '-'
	oldName := strings.TrimSpace(name)
	newNameSb := strings.Builder{}
	runes := []rune(oldName)
	dot := lastIndex(runes, '.')
	nameIndex := dot
	if dot < 0 {
		nameIndex = len(runes)
	}
	for i := 0; i < nameIndex; i++ {
		if mapping[runes[i]] {
			newNameSb.WriteString(midlineS)
		} else {
			newNameSb.WriteRune(runes[i])
		}
	}
	if dot >= 0 {
		for i := dot; i < len(runes); i++ {
			newNameSb.WriteRune(runes[i])
		}
	}

	r, _ := regexp.Compile("\\-+")
	newName := r.ReplaceAllString(newNameSb.String(), midlineS)
	newNameRunes := []rune(newName)
	if newNameRunes[0] == midlineR {
		newNameRunes = newNameRunes[1:]
	}
	newDotIndex := lastIndex(newNameRunes, '.')
	if newDotIndex <= 0 {
		if newNameRunes[len(newNameRunes)-1] == midlineR {
			return string(newNameRunes[0 : len(newNameRunes)-1])
		} else {
			return string(newNameRunes)
		}
	} else {
		if newNameRunes[newDotIndex-1] == midlineR {
			return string(newNameRunes[0:newDotIndex-1]) + string(newNameRunes[newDotIndex:])
		} else {
			return string(newNameRunes)
		}
	}
}

func lastIndex(runes []rune, r rune) int {
	dot := -1
	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] == r {
			dot = i
			break
		}
	}
	return dot
}
