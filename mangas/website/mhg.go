package website

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/zoroqi/rubbish/mangas/data"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Mhg struct {
}

const MhgSiteName = "mhg"

func (mhg *Mhg) ListParse(bs []byte) ([]data.Manga, error) {
	doc := soup.HTMLParse(string(bs))
	mangalist := doc.FindAll("div", "class", "latest-list")
	var manga []data.Manga
	for _, m2 := range mangalist {
		mm := m2.FindAll("li")
		for _, m := range mm {
			a := m.Find("p").Find("a")
			cover := m.Find("a", "class", "cover").Find("img")
			tt := m.Find("span", "class", "tt")
			fd := m.Find("span", "class", "fd")
			sl := m.Find("span", "class", "sl")
			dt := m.Find("span", "class", "dt")
			mangaType := ""

			if fd.Error != nil {
				mangaType += "1"
			} else if sl.Error != nil {
				mangaType += "2"
			}
			lastTT := ""
			if tt.Error == nil {
				lastTT = tt.FullText()
			}
			coverImg := ""
			if cover.Error == nil {
				coverImg = cover.Attrs()["data-src"]
				if coverImg == "" {
					coverImg = cover.Attrs()["src"]
				}
			}
			lastDateStr := ""
			if dt.Error == nil {
				lastDateStr = dt.FullText()
			}
			lastTT = strings.Replace(lastTT, "更新至", "", -1)

			lastDate, err := time.Parse(time.DateOnly, lastDateStr)
			if err != nil {
				fmt.Println(err, lastDateStr)
				lastDate = time.Now()
			}
			mg := data.Manga{
				Source:     MhgSiteName,
				Title:      strings.TrimSpace(a.Attrs()["title"]),
				Link:       strings.TrimSpace(a.Attrs()["href"]),
				Cover:      strings.TrimSpace(coverImg),
				Type:       strings.TrimSpace(mangaType),
				Last:       strings.TrimSpace(lastTT),
				LastUpdate: lastDate,
			}
			mg.Mid = ParseMhgMid(mg.Link)
			manga = append(manga, mg)
		}
	}
	return manga, nil
}

func (mhg *Mhg) ListParse2(bs []byte) ([]data.Manga, error) {
	doc := soup.HTMLParse(string(bs))
	mangalist := doc.FindAll("div", "class", "book-list")
	var manga []data.Manga
	for _, m2 := range mangalist {
		mm := m2.FindAll("li")
		for _, m := range mm {
			a := m.Find("p").Find("a")
			cover := m.Find("a", "class", "bcover").Find("img")
			tt := m.Find("span", "class", "tt")
			fd := m.Find("span", "class", "fd")
			sl := m.Find("span", "class", "sl")
			dt := m.Find("span", "class", "updateon")
			mangaType := ""

			if fd.Error != nil {
				mangaType += "1"
			} else if sl.Error != nil {
				mangaType += "2"
			}
			lastTT := ""
			if tt.Error == nil {
				lastTT = tt.FullText()
			}
			coverImg := ""
			if cover.Error == nil {
				coverImg = cover.Attrs()["data-src"]
				if coverImg == "" {
					coverImg = cover.Attrs()["src"]
				}
			}
			lastTT = strings.Replace(lastTT, "更新至", "", -1)

			lastDateStr := ""
			if dt.Error == nil {
				lastDateStr = strings.Replace(dt.Text(), "更新于：", "", -1)
				lastDateStr = strings.TrimSpace(lastDateStr)
			}
			lastDate, err := time.Parse(time.DateOnly, lastDateStr)
			if err != nil {
				lastDate = time.Now()
			}
			mg := data.Manga{
				Source:     MhgSiteName,
				Title:      strings.TrimSpace(a.Attrs()["title"]),
				Link:       strings.TrimSpace(a.Attrs()["href"]),
				Cover:      strings.TrimSpace(coverImg),
				Type:       strings.TrimSpace(mangaType),
				Last:       strings.TrimSpace(lastTT),
				LastUpdate: lastDate,
			}
			mg.Mid = ParseMhgMid(mg.Link)
			manga = append(manga, mg)
		}
	}
	return manga, nil
}

func (mhg *Mhg) LastUpdateList(link string) ([]byte, error) {
	return mhg.download(link, "https://www.manhuagui.com/")
}

func (mhg *Mhg) HistoryList(linkLayout string, page string, beforePage string) ([]byte, error) {
	return mhg.download(fmt.Sprintf(linkLayout, page), fmt.Sprintf(linkLayout, beforePage))
}
func (mhg *Mhg) download(link, referer string) ([]byte, error) {
	j, _ := url.Parse(link)
	req := &http.Request{}
	req.URL = j
	req.Header = make(http.Header)
	req.Header.Set("Referer", referer)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Method = http.MethodGet
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	if resp, err := client.Do(req); err != nil {
		return nil, fmt.Errorf("reqest erro, %s", err)
	} else {
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("status code %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("status code %d, read body err %w", resp.StatusCode, err)
		}
		return b, nil
	}
}

func (mhg *Mhg) Manga(link string) ([]byte, error) {
	return mhg.download(link, "https://www.manhuagui.com/")
}

func (mhg *Mhg) ParseManga(bs []byte) (data.MangaExtra, error) {
	doc := soup.HTMLParse(string(bs))
	bookCont := doc.Find("div", "class", "book-cont")
	detailList := bookCont.Find("ul", "class", "detail-list").FindAll("li")
	extra := data.MangaExtra{}
	if len(detailList) >= 1 {
		row1 := detailList[0].FindAll("span")
		if len(row1) >= 1 {
			c := row1[0].Find("a")
			extra.Create, _ = time.Parse("2006年", strings.TrimSpace(c.Text()))
		}
		if len(row1) >= 2 {
			c := row1[1].Find("a")
			extra.Country = strings.TrimSpace(c.Text())
		}
	}
	if len(detailList) >= 2 {
		row2 := detailList[1].FindAll("span")
		if len(row2) >= 1 {
			c := row2[0].FindAll("a")
			for _, cc := range c {
				extra.Type = append(extra.Type, strings.TrimSpace(cc.Text()))
			}
		}
	}
	if len(detailList) >= 3 {
		row3 := detailList[2].Find("span").FindAll("a")
		for _, r := range row3 {
			extra.Aliases = append(extra.Aliases, strings.TrimSpace(r.Text()))
		}
	}

	cl, err := mhg.ParseChapter(bs)
	if err != nil {
		return extra, err
	}
	extra.ChapterList = cl
	return extra, nil
}

func (mhg *Mhg) ParseChapter(bs []byte) ([]data.Chapter, error) {
	doc := soup.HTMLParse(string(bs))
	chapterList := doc.FindAll("div", "class", "chapter-list")
	var chapters []data.Chapter
	for _, cs := range chapterList {
		cl := cs.FindAll("li")
		for _, c := range cl {
			a := c.Find("a")
			title, chName := titleAndChapter(a.Attrs()["title"])
			ch := data.Chapter{
				Title:   title,
				Chapter: chName,
				Link:    strings.TrimSpace(a.Attrs()["href"]),
			}
			pageSize := a.Find("i")
			if pageSize.Error != nil {
				ch.Page = 0
			} else {
				ch.Page, _ = strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(pageSize.Text(), "p", "")))
			}
			ch.Type = chapterType(ch.Chapter)
			chapters = append(chapters, ch)
		}
	}
	return chapters, nil
}

func titleAndChapter(s string) (string, string) {
	ss := strings.SplitN(s, " ", 2)
	if len(ss) < 2 {
		return "", ss[0]
	}
	return ss[1], ss[0]
}

func chapterType(s string) data.ChapterType {
	if strings.Contains(s, "话") ||
		strings.Contains(s, "回") ||
		strings.Contains(s, "章") {
		return data.DanHua
	} else if strings.Contains(s, "卷") {
		return data.Juan
	} else if strings.Contains(s, "短篇") {
		return data.DuanPian
	} else if strings.Contains(s, "番外") {
		return data.FanWai
	} else if strings.Contains(s, "画集") {
		return data.HuaJi
	} else if strings.Contains(s, "合计") {
		return data.HeJi
	} else if strings.Contains(s, "特") {
		return data.TeBie
	} else if strings.Contains(s, "联动") {
		return data.LianDong
	} else {
		return data.Other
	}
}

// https://www.manhuagui.com/comic/16460/
// https://www.manhuagui.com/comic/28004/
func ParseMhgMid(link string) string {
	p := regexp.MustCompile(`comic/(\d+)/?`)
	r := p.FindStringSubmatch(link)
	if len(r) <= 1 {
		return link
	}
	return r[1]
}
