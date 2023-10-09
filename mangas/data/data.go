package data

import (
	"strconv"
	"strings"
	"time"
)

type Manga struct {
	Mid        string    `json:"mid"`
	Source     string    `json:"source"`
	Title      string    `json:"title"`
	Link       string    `json:"link"`
	Type       string    `json:"type"`
	Cover      string    `json:"cover"`
	Last       string    `json:"last"`
	LastUpdate time.Time `json:"lastUpdate"`
}

func ParseSlice(s []string) Manga {
	r := Manga{}
	r.Source = s[0]
	r.Title = s[1]
	r.Link = s[2]
	r.Cover = s[3]
	r.Type = s[4]
	r.Last = s[5]

	if len(s) >= 7 {
		if s[6] != "" {
			r.LastUpdate, _ = time.Parse(time.DateOnly, s[6])
		}
	}

	if len(s) >= 8 {
		r.Mid = s[7]
	}
	return r
}

func (m Manga) ToSlice() []string {
	r := []string{
		m.Source,
		m.Title,
		m.Link,
		m.Cover,
		m.Type,
		m.Last,
		m.LastUpdate.Format(time.DateOnly),
		m.Mid,
	}
	return r
}

type MangaExtra struct {
	Manga
	Country     string    `json:"country"`
	Authors     []string  `json:"authors"`
	Aliases     []string  `json:"aliases"`
	Type        []string  `json:"type"`
	Create      time.Time `json:"create"`
	ChapterList []Chapter `json:"chapterList"`
}

func (m MangaExtra) ToSlice() []string {
	r := m.Manga.ToSlice()
	r = append(r, m.Country)
	r = append(r, strings.Join(m.Authors, ","))
	r = append(r, strings.Join(m.Aliases, ","))
	r = append(r, strings.Join(m.Type, ","))
	r = append(r, m.Create.Format(time.DateOnly))
	return r
}

type Chapter struct {
	Mid     string      `json:"mid"`
	Chapter string      `json:"chapter"`
	Link    string      `json:"link"`
	Title   string      `json:"title"`
	Page    int         `json:"page"`
	Type    ChapterType `json:"type"`
}

func (c Chapter) ToSlice() []string {
	r := []string{
		c.Mid,
		c.Chapter,
		c.Link,
		c.Title,
		strconv.Itoa(c.Page),
		strconv.Itoa(int(c.Type)),
	}
	return r
}

type ChapterType int

const (
	Other = iota
	DanHua
	Juan
	DuanPian
	HeJi
	TeBie
	FanWai
	HuaJi
	LianDong
)
