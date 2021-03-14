package main

import (
	"flag"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	//"strings"
)

const (
	explore  = "explore"
	trending = "trending"
)

type exploreConfig struct {
	Name  string `yaml:"name"`
	Param string `yaml:"param"`
}

// https://github.com/explore
func executeExplore() (string, error) {
	return "", nil
}

// https://github.com/trending/c++?since=weekly
func executeTrending(config exploreConfig) (string, error) {

	f, err := httpget("https://github.com/trending/" + config.Param)
	if err != nil {
		return "", err
	}
	root := soup.HTMLParse(f)
	rows := root.FindAll("article", "class", "Box-row")
	text := ""
	for i, r := range rows {
		repo := r.Find("h1").FullText()
		p := r.Find("p")
		desc := ""
		if p.Error == nil {
			desc = strings.TrimSpace(p.FullText())
		}
		r := regexp.MustCompile("\\s+")
		repo = r.ReplaceAllString(repo, "")

		text += fmt.Sprintf("%d. [%s](https://github.com/%s) %s\n", i, repo, repo, desc)
	}
	return text, nil
}

func httpget(url string) (string, error) {
	req := &fasthttp.Request{}
	req.SetRequestURI(url)
	req.Header.SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.100 Safari/537.36")
	req.Header.SetMethod("GET")
	resp := &fasthttp.Response{}
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return "", err
	}
	defer resp.ConnectionClose()
	b := resp.Body()
	return string(b), nil
}

func mkdir(dirpath string) error {
	return os.MkdirAll(dirpath, 0755)
}

func main() {

	config := flag.String("c", "", "config path")
	output := flag.String("o", "", "output directory")
	collcetExplore := flag.Bool("e", false, "collect explore")
	flag.Parse()
	if *config == "" || *output == "" {
		fmt.Println("'c' or 'o' is empty")
		return
	}

	var configs []exploreConfig

	configStr, err := ioutil.ReadFile(*config)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = yaml.Unmarshal(configStr, &configs)
	if err != nil {
		fmt.Println(err)
		return
	}

	t := time.Now()
	year := strconv.Itoa(t.Year())
	month := strconv.Itoa(int(t.Month()))
	date := t.Format("2006-01-02") + ".md"

	trendingOutput := filepath.Join(*output, trending, year, month)
	if err := mkdir(trendingOutput); err != nil {
		return
	}

	contents := `# ` + t.Format("2006-01-02") + "\n"
	text := ""
	for i, c := range configs {
		f, err := executeTrending(c)
		if err != nil {
			continue
		}
		fmt.Println(i, c.Name)
		contents += fmt.Sprintf("* [%s](#%s)\n", c.Name, c.Name)
		text += fmt.Sprintf("# %s\n\n%s\n\n", c.Name, f)

	}

	trendingText := contents + "\n\n" + text
	fmt.Println("output trending ", ioutil.WriteFile(filepath.Join(trendingOutput, date), []byte(trendingText), 0666))

	if *collcetExplore {
		exploreOutput := filepath.Join(*output, explore, year, month, date)
		f, err := executeExplore()
		if err != nil {
			return
		}
		fmt.Println("output explore ", ioutil.WriteFile(filepath.Join(exploreOutput, date), []byte(f), 0666))
	}
}
