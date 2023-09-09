package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/google/go-github/v55/github"
	"github.com/spf13/cobra"
	"github.com/zoroqi/rubbish/mangas/data"
	"github.com/zoroqi/rubbish/mangas/website"
	"golang.org/x/oauth2"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	githubtoken = os.Getenv("GITHUB_TOKEN")
	repos       string
	user        string
	outfile     string
)

var updateScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan manga website",
	Long:  "scan manga website",
	RunE:  updatelist,
}

func init() {
	RootCmd.AddCommand(updateScanCmd)
	updateScanCmd.Flags().StringVar(&repos, "repos", "", "github repos")
	updateScanCmd.Flags().StringVar(&user, "user", "", "github user")
	updateScanCmd.Flags().StringVar(&outfile, "f", "", "output file")
}

func githubCommit(web, year, month, fileName string, content string) error {
	if githubtoken == "" || repos == "" || user == "" {
		return nil
	}

	owner := user
	repo := repos

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

func updatelist(cmd *cobra.Command, args []string) error {
	rand.Seed(time.Now().UnixNano())
	yesterday := time.Now().AddDate(0, 0, -1)
	if outfile == "" {
		outfile = yesterday.Format(time.DateOnly) + ".csv"
		fmt.Println("f empty, default out file name:", outfile)
	}

	mhg := &website.Mhg{}
	mgs, err := ListUpdate(mhg, "https://www.manhuagui.com/update/")
	if err != nil {
		fmt.Println(err)
		return err
	}
	file := &strings.Builder{}
	writeFile(file, saveToday(mgs, yesterday.Format(time.DateOnly)))
	fmt.Println(file.String())
	if err := githubCommit(website.MhgSiteName, yesterday.Format("2006"), yesterday.Format("01"),
		outfile, file.String()); err != nil {
		fmt.Println("github commit err", err)
		return err
	}
	return nil
}

func saveToday(mgs []data.Manga, date string) []data.Manga {
	r := []data.Manga{}
	for _, m := range mgs {
		if m.LastUpdate.Format(time.DateOnly) == date {
			r = append(r, m)
		}
	}
	return r
}

func ListUpdate(mhg *website.Mhg, url string) ([]data.Manga, error) {
	bs, err := mhg.LastUpdateList(url)
	if err != nil {
		return nil, err
	}
	return mhg.ListParse(bs)
}

func writeFile(file io.Writer, listResult []data.Manga) {
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
