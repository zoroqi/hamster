package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	var infos []info
	err := filepath.Walk(root, func(path string, fileInfo os.FileInfo, err error) error {
		if err == nil && fileInfo.IsDir() && !strings.HasPrefix(fileInfo.Name(), ".") {
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
	replaceName(infos)
}

func replaceName(infos []info) {
	mapping := loadMapping()
	if mapping == nil {
		return
	}
	for _, info := range infos {
		newName := rename(info.fileInfo.Name(), mapping)
		if newName != info.fileInfo.Name() {
			fmt.Println(info.fileInfo.Name() + " -> " + newName)
			path := info.path
			index := strings.LastIndex(path, "/")
			var err error
			if index < 0 {
				err = os.Rename(info.path, newName)
			} else {
				newPath := path[0:index] + "/" + newName
				err = os.Rename(info.path, newPath)
			}
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
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
	replaceName(infos)
}

func loadMapping() map[rune]bool {
	return defaultMapping()
}

func defaultMapping() map[rune]bool {
	mapping := make(map[rune]bool)
	mapping[' '] = true
	mapping['　'] = true

	mapping['：'] = true
	mapping[':'] = true

	mapping['；'] = true
	mapping[';'] = true

	mapping['，'] = true
	mapping['。'] = true
	mapping['、'] = true

	mapping['《'] = true
	mapping['》'] = true

	mapping['—'] = true

	mapping['（'] = true
	mapping['）'] = true

	mapping['+'] = true

	mapping['【'] = true
	mapping['】'] = true
	mapping['|'] = true
	mapping['！'] = true
	mapping['!'] = true
	mapping['／'] = true
	mapping['\\'] = true
	mapping['〜'] = true
	mapping['~'] = true
	mapping['\''] = true
	mapping['"'] = true

	mapping['’'] = true
	mapping['‘'] = true
	mapping['”'] = true
	mapping['“'] = true
	mapping['…'] = true
	mapping['*'] = true
	mapping[','] = true
	mapping[')'] = true
	mapping['('] = true
	mapping['['] = true
	mapping[']'] = true
	mapping['&'] = true
	mapping['#'] = true
	mapping['^'] = true
	mapping['`'] = true
	mapping['@'] = true
	mapping['・'] = true
	mapping['{'] = true
	mapping['}'] = true
	mapping['「'] = true
	mapping['」'] = true

	return mapping
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
