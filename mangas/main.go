package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/google/go-github/v55/github"
	"golang.org/x/oauth2"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Manga struct {
	Source     string    `json:"source"`
	Title      string    `json:"title"`
	Link       string    `json:"link"`
	Type       string    `json:"type"`
	Cover      string    `json:"cover"`
	Last       string    `json:"last"`
	LastUpdate time.Time `json:"lastUpdate"`
}

var (
	githubtoken = os.Getenv("GITHUB_TOKEN")
	repos       = flag.String("repos", "", "github repos")
	user        = flag.String("user", "", "github user")
)

func init() {
	flag.Parse()
}

func githubCommit(web, year, month, fileName string, content string) error {
	if githubtoken == "" || *repos == "" || *user == "" {
		return nil
	}

	owner := *user
	repo := *repos

	client := newClient(githubtoken)
	ref, _, err := client.Git.GetRef(context.Background(),
		owner, repo, fmt.Sprintf("heads/main"))
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/%s/%s/%s", web, year, month, fileName)

	var entries []*github.TreeEntry
	entries = append(entries, &github.TreeEntry{Path: github.String(path),
		Type:    github.String("blob"),
		Content: github.String(content),
		Mode:    github.String("100644"),
	})

	tree, _, err := client.Git.CreateTree(context.Background(), owner, repo, *ref.Object.SHA, entries)
	if err != nil {
		return err
	}

	parent, _, err := client.Repositories.GetCommit(context.Background(), owner, repo, *ref.Object.SHA, nil)
	if err != nil {
		return err
	}
	parent.Commit.SHA = parent.SHA
	message := "commit " + fileName
	commitName := "github-actions[bot]"
	commitEmail := "41898282+github-actions[bot]@users.noreply.github.com"
	date := github.Timestamp{Time: time.Now()}
	author := &github.CommitAuthor{Date: &date, Name: &commitName, Email: &commitEmail}
	commit := &github.Commit{Author: author,
		Message: &message,
		Tree:    tree,
		Parents: []*github.Commit{parent.Commit}}
	newCommit, _, err := client.Git.CreateCommit(context.Background(), owner, repo, commit)
	if err != nil {
		return err
	}
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(context.Background(), owner, repo, ref, false)
	return err
}

func newClient(token string) *github.Client {
	var client *github.Client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client = github.NewClient(tc)
	return client
}

func main() {
	rand.Seed(time.Now().UnixNano())
	f := flag.String("f", "", "output file")
	flag.Parse()
	yesterday := time.Now().AddDate(0, 0, -1)
	if *f == "" {
		*f = yesterday.Format(time.DateOnly) + ".csv"
		fmt.Println("f empty, default out file name:", *f)
	}

	mhg := &Mhg{}
	mgs, err := ListUpdate(mhg, "https://www.manhuagui.com/update/")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	file := &strings.Builder{}
	writeFile(file, saveToday(mgs, yesterday.Format(time.DateOnly)))
	fmt.Println(file.String())
	if err := githubCommit(MhgSiteName, yesterday.Format("2006"), yesterday.Format("01"),
		*f, file.String()); err != nil {
		fmt.Println("github commit err", err)
		os.Exit(1)
	}
}

func saveToday(mgs []Manga, date string) []Manga {
	r := []Manga{}
	for _, m := range mgs {
		if m.LastUpdate.Format(time.DateOnly) == date {
			r = append(r, m)
		}
	}
	return r
}

func ListUpdate(mhg *Mhg, url string) ([]Manga, error) {
	bs, err := mhg.LastUpdateList(url)
	if err != nil {
		return nil, err
	}
	return mhg.ListParse(bs)
}

func writeFile(file io.Writer, listResult []Manga) {
	writer := csv.NewWriter(file)
	count := 0
	writer.Write([]string{"source", "title", "link", "cover", "type", "last"})
	for _, manga := range listResult {
		err := writer.Write([]string{manga.Source, manga.Title, manga.Link, manga.Cover, manga.Type, manga.Last})
		if err != nil {
			fmt.Printf("marshal err, %+v, %v\n", manga, err)
			continue
		}
		count++
	}
	writer.Flush()
	fmt.Printf("write end, write:%d\n", count)
	if e := writer.Error(); e != nil {
		fmt.Println(e)
		return
	}
}
