package manga_search

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/zoroqi/rubbish/manga"
	"io"
	"os"
	"strings"
)

type memoryMangaDb struct {
	data  []manga.Manga
	index map[rune][]mangaIndex
}

func (m *memoryMangaDb) Search(str string) []manga.Manga {
	if str == "" {
		return nil
	}
	runes := []rune(str)

	ms := m.index[runes[0]]

	if len(ms) == 0 {
		return nil
	}
	r := make([]manga.Manga, 0)
	for _, i := range ms {
		manga := m.data[i.index]
		if strings.Contains(manga.Title, str) {
			r = append(r, manga)
		} else {
			for _, t := range manga.Titles {
				if strings.Contains(t.Name, str) {
					r = append(r, manga)
				}
			}
		}

	}
	return r
}

type dustman []removerHandler

func (d dustman) remove(manga manga.Manga) bool {
	for _, h := range d {
		if h(manga) {
			return true
		}
	}
	return false
}

type removerHandler func(m manga.Manga) bool

type dupRemover map[string]bool

func (d dupRemover) remover(m manga.Manga) bool {
	if d[m.Id] {
		return true
	} else {
		d[m.Id] = true
		return false
	}
}

type mangaIndex struct {
	runes []rune
	index int
}

func newMangIndex(index int, m manga.Manga) []mangaIndex {
	r := make([]mangaIndex, 0, 3)
	r = append(r, mangaIndex{runes: []rune(m.Title), index: index})
	for _, t := range m.Titles {
		r = append(r, mangaIndex{runes: []rune(t.Name), index: index})
	}
	return r
}

func Load(file string) (*memoryMangaDb, error) {

	var removers dustman
	removers = append(removers, dupRemover{}.remover)

	mangas, err := readFile(file, removers)
	if err != nil {
		return nil, err
	}

	db := &memoryMangaDb{}
	db.data = mangas
	index := make(map[rune][]mangaIndex)
	for i, mg := range mangas {
		mi := newMangIndex(i, mg)
		for _, m := range mi {
			for _, r := range m.runes {
				if a, ok := index[r]; ok {
					index[r] = append(a, m)
				} else {
					index[r] = append([]mangaIndex{}, m)
				}
			}
		}
	}
	db.index = index

	return db, nil
}

func readFile(file string, dust dustman) ([]manga.Manga, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var initLength int64
	if info, err := f.Stat(); err == nil {
		initLength = info.Size() / 100
	} else {
		initLength = 10000
	}
	read := bufio.NewReader(f)
	mangas := make([]manga.Manga, 0, initLength)
	for {
		line, err := read.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Read file error!", err)
				return nil, err
			}
		}
		if len(line) == 0 {
			continue
		}
		var m manga.Manga
		if err = json.Unmarshal(line, &m); err != nil {
			fmt.Printf("json parse err, %s, %+v", string(line), err)
			continue
		}
		if dust.remove(m) {
			continue
		}
		mangas = append(mangas, m)
	}
	return mangas, nil
}
