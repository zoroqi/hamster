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
		return err
	})
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	replaceName(infos)
}

func replaceName(infos []info) {
	mapping := loadMapping("")
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
		return err
	})
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	replaceName(infos)
}

func loadMapping(path string) map[rune]string {
	if path == "" {
		return defaultMapping()
	}
	return nil
}

func defaultMapping() map[rune]string {
	mapping := make(map[rune]string)
	mapping[' '] = "-"
	mapping['　'] = "-"

	mapping['：'] = "-"
	mapping[':'] = "-"

	mapping['；'] = "-"
	mapping[';'] = "-"

	mapping['，'] = "-"
	mapping['。'] = "-"
	mapping['、'] = "-"

	mapping['《'] = "-"
	mapping['》'] = "-"

	mapping['—'] = "-"

	mapping['（'] = "-"
	mapping['）'] = "-"

	mapping['+'] = "-"

	mapping['【'] = "-"
	mapping['】'] = "-"
	mapping['|'] = "-"
	mapping['！'] = "-"
	mapping['!'] = "-"
	mapping['／'] = "-"
	mapping['\\'] = "-"
	mapping['〜'] = "-"
	mapping['~'] = "-"
	mapping['\''] = "-"
	mapping['"'] = "-"

	mapping['’'] = "-"
	mapping['‘'] = "-"
	mapping['”'] = "-"
	mapping['“'] = "-"
	mapping['…'] = "-"
	mapping['*'] = "-"
	mapping[','] = "-"
	mapping[')'] = "-"
	mapping['('] = "-"
	mapping['['] = "-"
	mapping[']'] = "-"
	mapping['&'] = "-"
	mapping['#'] = "-"
	mapping['^'] = "-"
	mapping['`'] = "-"
	mapping['@'] = "-"
	mapping['・'] = "-"
	mapping['{'] = "-"
	mapping['}'] = "-"
	mapping['「'] = "-"
	mapping['」'] = "-"

	return mapping
}

func rename(name string, mapping map[rune]string) string {
	oldName := strings.TrimSpace(name)
	newNameSb := strings.Builder{}
	runes := []rune(oldName)
	dot := lastIndex(runes, '.')
	noExt := dot < 0
	if noExt {
		dot = len(runes)
	}
	for i := 0; i < dot; i++ {
		if m, exist := mapping[runes[i]]; exist {
			newNameSb.WriteString(m)
		} else {
			newNameSb.WriteRune(runes[i])
		}
	}
	if !noExt {
		for i := dot; i < len(runes); i++ {
			newNameSb.WriteRune(runes[i])
		}
	}

	r, _ := regexp.Compile("\\-+")
	newName := r.ReplaceAllString(newNameSb.String(), "-")

	newNameRunes := []rune(newName)
	if newNameRunes[0] == '-' {
		newNameRunes = newNameRunes[1:]
	}
	newDotIndex := lastIndex(newNameRunes, '.')
	if newDotIndex <= 0 {
		return string(newNameRunes)
	} else {
		if newNameRunes[newDotIndex-1] == '-' {
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
