package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/zoroqi/hamster/notes-analysis/db"
	"github.com/zoroqi/hamster/notes-analysis/parser"
	"github.com/zoroqi/hamster/notes-analysis/store"
	"os"
	"strings"
)

var (
	loadflag = flag.Bool("l", false, "load")
	loadfile = flag.String("lf", "", "load file")
	dumpfile = flag.String("df", "", "dump file")
	scanPath = flag.String("scan", "", "scan path")
)

func main() {

	flag.Parse()
	var docs []*store.Document
	var err error
	if *loadflag {
		if *loadfile == "" {
			fmt.Println("no load file")
			return
		}
		docs, err = load(*loadfile)
	} else {
		if *dumpfile == "" {
			fmt.Println("no dump file")
			return
		}
		err = dump(*dumpfile)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	tags, note2tag := hashtags(docs)
	client, err := db.NewClient()
	defer client.Close(context.TODO())
	if err != nil {
		return
	}
	db.SaveTagsDQL(client, tags)

	nodes := parser.MakingGraph(docs)
	walk(client, nodes, note2tag)
}

//err = parser.WalkGraph(nodes, func(from *store.GraphNode, to *store.GraphNode) error {
//	//if to == nil {
//	//	fmt.Println(from.Doc.Name, " -> not found")
//	//	return nil
//	//}
//	if to != nil {
//		fmt.Print(from.Id, " -> ", to.Id)
//	} else {
//		fmt.Print(from.Id, " -> no")
//	}
//	if from.Block != nil {
//		fmt.Print(from.Block.HashTags)
//	}
//	if to != nil && to.Block != nil {
//		fmt.Print(to.Block.HashTags)
//	}
//	fmt.Println("")
//	return nil
//})
//}

func charsetEscapes(s string) string {
	r := []rune(s)
	for i := range r {
		if r[i] == '"' {
			r[i] = '\''
		}
	}
	return string(r)
}

func hashtags(nodes []*store.Document) (map[string]*db.DGraphTag, map[string][]*db.DGraphTag) {
	hashtags := map[string]*db.DGraphTag{}

	multiTag2 := func(ss []string) *db.DGraphTag {
		var r *db.DGraphTag
		for i := range ss {
			tt := strings.Join(ss[:i+1], "/")
			if t, ok := hashtags[tt]; !ok {
				t := &db.DGraphTag{
					Id:    fmt.Sprintf("%x", md5.Sum([]byte(tt)))[:8],
					Tag:   tt,
					Child: make(map[string]*db.DGraphTag),
				}
				hashtags[strings.Join(ss[:i+1], "/")] = t
				r = t
			} else {
				r = t
			}
		}
		for i := range ss {
			if i == 0 {
				continue
			}
			pt := strings.Join(ss[:i], "/")
			tt := strings.Join(ss[:i+1], "/")
			hashtags[pt].Child[tt] = hashtags[tt]
		}
		return r
	}
	multiTag := func(t string) *db.DGraphTag {
		ss := strings.Split(t, "/")
		return multiTag2(ss)
	}

	node2Tag := map[string][]*db.DGraphTag{}

	for _, doc := range nodes {
		meta := doc.Meta
		if meta != nil {
			for _, v := range meta {
				if v.Key == "tags" {
					if tags, ok := v.Value.([]any); ok {
						for _, tag := range tags {
							if t, ok := tag.(string); ok {
								node2Tag[doc.Id] = append(node2Tag[doc.Id], multiTag(t))
							}
						}
					}
				}
			}
		}

		for _, block := range doc.Blocks {
			store.WalkBlocks(block, func(b *store.Block) bool {
				for _, ht := range b.HashTags {
					node2Tag[b.Id] = append(node2Tag[b.Id], multiTag2(ht))
				}
				return true
			})
		}
	}
	return hashtags, node2Tag
}
func walk(ctx neo4j.DriverWithContext, nodes map[string]*store.GraphNode, note2tag map[string][]*db.DGraphTag) error {
	dup := map[string]bool{}
	err := parser.WalkGraph(nodes, func(from *store.GraphNode, to *store.GraphNode) error {
		if from.Doc != nil && !dup[from.Doc.Id] {
			dup[from.Doc.Id] = true
			db.PrintDoc(ctx, from.Doc, note2tag)
		}
		//if from.Block != nil && !dup[from.Block.Id] {
		//	dup[from.Block.Id] = true
		//	db.PrintBlock(ctx, from.Block, note2tag)
		//}
		//if to != nil {
		//	db.PrintDep(ctx, from, to)
		//}
		return nil
	})
	err = parser.WalkGraph(nodes, func(from *store.GraphNode, to *store.GraphNode) error {
		if to != nil {
			if from.Id != to.Id {
				db.PrintDep(ctx, from, to)
			}
		}
		return nil
	})
	return err
}

func dump(dumpfile string) error {
	filePath := *scanPath
	if filePath == "" {
		return errors.New("no scan path")
	}
	docs, err := parser.ParseAllfile(filePath)
	if err != nil {
		return err
	}

	for _, doc := range docs {
		for _, b := range doc.Blocks {
			store.WalkBlocks(b, func(b *store.Block) bool {
				r := []rune(b.Content)
				b.Content = string(r[:min(20, len(r))])
				return true
			})
		}
	}

	bs, err := json.MarshalIndent(docs, "", " ")
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return os.WriteFile(dumpfile, bs, 0644)
}

func load(loadfile string) ([]*store.Document, error) {
	bs, err := os.ReadFile(loadfile)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	docs := []*store.Document{}
	err = json.Unmarshal(bs, &docs)
	return docs, nil
}
