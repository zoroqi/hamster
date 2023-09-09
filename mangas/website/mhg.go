package website

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/zoroqi/rubbish/mangas/data"
	"io"
	"net/http"
	"net/url"
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
			manga = append(manga, data.Manga{
				Source:     MhgSiteName,
				Title:      strings.TrimSpace(a.Attrs()["title"]),
				Link:       strings.TrimSpace(a.Attrs()["href"]),
				Cover:      strings.TrimSpace(coverImg),
				Type:       strings.TrimSpace(mangaType),
				Last:       strings.TrimSpace(lastTT),
				LastUpdate: lastDate,
			})
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
				lastDateStr = strings.Replace(dt.FullText(), "更新于：", "", -1)
				lastDateStr = strings.TrimSpace(lastDateStr)
			}
			lastDate, err := time.Parse(time.DateOnly, lastDateStr)
			if err != nil {
				lastDate = time.Now()
			}
			manga = append(manga, data.Manga{
				Source:     MhgSiteName,
				Title:      strings.TrimSpace(a.Attrs()["title"]),
				Link:       strings.TrimSpace(a.Attrs()["href"]),
				Cover:      strings.TrimSpace(coverImg),
				Type:       strings.TrimSpace(mangaType),
				Last:       strings.TrimSpace(lastTT),
				LastUpdate: lastDate,
			})
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
