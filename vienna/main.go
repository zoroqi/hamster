package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"github.com/zoroqi/rubbish/clipboard"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	db     = flag.String("db", "", "vienna db path")
	output = flag.String("output", "", "output file path")
)

func main() {
	flag.Parse()
	var sqlite3conn *sqlite3.SQLiteConn
	sql.Register("messages", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			sqlite3conn = conn
			return nil
		},
	})
	if *db == "" || *output == "" {
		fmt.Println("db or output is black")
		return
	}
	db, err := sql.Open("messages", *db)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	rssFolderMapping, err := loadRss(db)
	if err != nil {
		fmt.Println(err)
		return
	}
	marks, err := loadMarked(db, rssFolderMapping)
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range marks {
		path := fmt.Sprintf("%s/%s.md", *output, k)
		f, err := os.Open(path)
		var old []clipboard.Mark
		if err == nil {
			old, _ = clipboard.ParseStoreFile(f)
			f.Close()
		}
		mark := clipboard.Merge(old, v)

		ioutil.WriteFile(path, []byte(clipboard.FileString(k, mark).String()), 0666)
	}
}

type rss struct {
	link  string
	title string
}

func loadRss(db *sql.DB) (map[int64]rss, error) {
	rows, err := db.Query("select folder_id,description,home_page from rss_folders ")
	if err != nil {
		return nil, err
	}
	result := make(map[int64]rss)
	for rows.Next() {
		var homePage string
		var description string
		var folderId int64
		err = rows.Scan(&folderId, &description, &homePage)
		if err != nil {
			return nil, err
		}
		m := rss{
			link:  strings.TrimSpace(homePage),
			title: strings.TrimSpace(description),
		}
		if len(m.title) == 0 {
			m.title = m.link
		}
		if len(m.link) == 0 {
			continue
		}
		result[folderId] = m
	}
	return result, nil
}

func loadMarked(db *sql.DB, rssFolderMapping map[int64]rss) (map[string][]clipboard.Mark, error) {
	rows, err := db.Query("select folder_id,title,link,date from messages where marked_flag=1 order by createddate desc")
	if err != nil {
		return nil, err
	}
	result := make(map[string][]clipboard.Mark)
	for rows.Next() {
		var title string
		var link string
		var date float64
		var folderId int64
		err = rows.Scan(&folderId, &title, &link, &date)
		if err != nil {
			fmt.Println(err)
			continue
		}
		createDate := time.Unix(int64(date), 0)
		m := clipboard.Mark{
			Text:    title,
			Link:    strings.TrimSpace(link),
			LinkStr: fmt.Sprintf("[%s][%s]", strings.TrimSpace(title), strings.TrimSpace(link)),
		}
		if rs, exist := rssFolderMapping[folderId]; exist {
			m.Source = fmt.Sprintf("[%s][%s]", rs.title, rs.link)
		}
		if uri, err := url.Parse(link); err == nil {
			r := regexp.MustCompile("([^.]*?)\\.([^.]*?)$")
			ss := r.FindAllStringSubmatch(uri.Hostname(), -1)
			if len(ss) == 1 {
				m.Tags = []string{ss[0][1]}
			}
		}
		m.Id = clipboard.NewId(m.Link)
		result[createDate.Format("200601")] = append(result[createDate.Format("200601")], m)
	}

	return result, nil
}
