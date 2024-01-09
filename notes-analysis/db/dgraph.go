package db

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/zoroqi/hamster/notes-analysis/store"
	"strings"
	"time"
)

type School struct {
	Name  string   `json:"name,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

type loc struct {
	Type   string    `json:"type,omitempty"`
	Coords []float64 `json:"coordinates,omitempty"`
}

// If omitempty is not set, then edges with empty values (0 for int/float, "" for string, false
// for bool) would be created for values not specified explicitly.

type Person struct {
	Uid      string     `json:"uid,omitempty"`
	Name     string     `json:"name,omitempty"`
	Age      int        `json:"age,omitempty"`
	Dob      *time.Time `json:"dob,omitempty"`
	Married  bool       `json:"married,omitempty"`
	Raw      []byte     `json:"raw_bytes,omitempty"`
	Friends  []Person   `json:"friend,omitempty"`
	Location loc        `json:"loc,omitempty"`
	School   []School   `json:"school,omitempty"`
	DType    []string   `json:"dgraph.type,omitempty"`
}

func NewClient() (neo4j.DriverWithContext, error) {
	ctx := context.Background()
	// URI examples: "neo4j://localhost", "neo4j+s://xxx.databases.neo4j.io"
	dbUri := "neo4j://localhost:7687"
	dbUser := "neo4j"
	dbPassword := "12345678"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))

	err = driver.VerifyConnectivity(ctx)

	return driver, err
}

// block, doc, tag
// doc -> block: file
// block -> doc: link
// doc -> doc: link
// doc,block -> tag
//

type DGraphBlockNode struct {
	Id      string
	Text    string
	Doc     *DGraphDocNode
	LinkDoc *DGraphDocNode
	Blocks  []DGraphBlockNode
}

type DGraphDocNode struct {
	Id     string
	Tags   []DGraphTag
	Blocks []DGraphBlockNode
}

type DGraphTag struct {
	Id    string
	Tag   string
	Child map[string]*DGraphTag
}

func printfln(format string, a ...any) (n int, err error) {
	return fmt.Printf(format+"\n", a...)
}

func cutId(s string) string {
	return s[0:min(len(s), 8)]
}

/*
alter 内容

type ADoc {
	id
	path
	name
	extra
	tags
	doctype
	mod_time
}

id: string @index(term) @lang .
level: int  .
path: string .
name: string .
doctype: string .
extra: string .
tags: [uid] .
mod_time: dateTime .
*/

func PrintDoc(ctx neo4j.DriverWithContext, doc *store.Document, note2tag map[string][]*DGraphTag) {
	id := cutId(doc.Id)
	name := doc.Name
	if doc.Meta != nil {
		for _, v := range doc.Meta {
			if v.Key == "aliases" {
				if t, ok := v.Value.([]any); ok {
					if len(t) > 0 {
						if n, ok := t[0].(string); ok {
							if strings.TrimSpace(n) != "" {
								name = strings.TrimSpace(n)
							}
						}
					}
				}
			}
		}
	}
	cql := fmt.Sprintf("CREATE (n: ADoc {id:$id, path:$path, name:$name, extra:$extra, doctype:$doctype, mod_time:$mod_time}) RETURN n;")
	result, err := neo4j.ExecuteQuery(context.Background(), ctx,
		cql,
		map[string]any{
			"id":       id,
			"path":     doc.Path,
			"name":     name,
			"extra":    doc.Extra,
			"doctype":  doc.Type,
			"mod_time": doc.ModTime.Format(time.RFC3339),
		}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))
	fmt.Println(result, err, cql, doc.Id, doc.Name)

	dup := map[string]bool{}
	for _, tag := range note2tag[doc.Id] {
		if !dup[tag.Id] {
			dup[tag.Id] = true
			cql := fmt.Sprintf("MATCH (a:ADoc {id:$did}), (b:HashTag {id:$tid}) MERGE (a)-[:TAGS]->(b) MERGE (b)-[:TAGS]->(a)")
			result, err := neo4j.ExecuteQuery(context.Background(), ctx,
				cql,
				map[string]any{
					"did": id,
					"tid": tag.Id,
				}, neo4j.EagerResultTransformer,
				neo4j.ExecuteQueryWithDatabase("neo4j"))
			fmt.Println(result, err, cql, doc.Id, tag.Id, tag.Tag)
		}
	}

	for _, block := range doc.Blocks {
		store.WalkBlocks(block, func(b *store.Block) bool {
			for _, tag := range note2tag[b.Id] {
				if !dup[tag.Id] {
					dup[tag.Id] = true
					cql := fmt.Sprintf("MATCH (a:ADoc {id:$did}), (b:HashTag {id:$tid}) MERGE (a)-[:TAGS]->(b) MERGE (b)-[:TAGS]->(a)")
					result, err := neo4j.ExecuteQuery(context.Background(), ctx,
						cql,
						map[string]any{
							"did": id,
							"tid": tag.Id,
						}, neo4j.EagerResultTransformer,
						neo4j.ExecuteQueryWithDatabase("neo4j"))
					fmt.Println(result, err, cql, doc.Id, tag.Id, tag.Tag)
				}
			}
			return true
		})
	}
}

func PrintDep(ctx neo4j.DriverWithContext, from *store.GraphNode, to *store.GraphNode) {
	fid := cutId(from.Id)
	tid := cutId(to.Id)
	cql := fmt.Sprintf("MATCH (a:ADoc {id:$fid}), (b:ADoc {id:$tid}) MERGE (a)-[:REF]->(b) MERGE (b)-[:REFD]->(a)")
	result, err := neo4j.ExecuteQuery(context.Background(), ctx,
		cql,
		map[string]any{
			"fid": fid,
			"tid": tid,
		}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))
	fmt.Println(result, err, cql, fid, tid)
}

/*
   type ABlock {
   	id
   	level
   	kind
   	content
   	tags
   }

   id: string @index(term) @lang .
   level: int  .
   kind: string .
   content: string .
   tags: [uid] .
*/

func PrintBlock(ctx neo4j.DriverWithContext, block *store.Block, note2tag map[string][]*DGraphTag) {
	id := cutId(block.Id)
	cql := fmt.Sprintf("CREATE (n: ABlock {id:$id, level:$level, kind:$kind, content:$content}) RETURN n;")
	result, err := neo4j.ExecuteQuery(context.Background(), ctx,
		cql,
		map[string]any{
			"id":      id,
			"level":   block.Level,
			"kind":    block.Kind,
			"content": block.Content,
		}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))
	fmt.Println(result, err, cql, block.Id, block.Content)
	store.WalkBlocks(block, func(b *store.Block) bool {
		dup := map[string]bool{}
		for _, tag := range note2tag[b.Id] {
			if !dup[tag.Id] {
				dup[tag.Id] = true
				cql := fmt.Sprintf("MATCH (a:ABlock {id:$bid}), (b:HashTag {id:$tid}) MERGE (a)-[:TAGS]->(b) MERGE  (b)-[:TAGS]->(a)")
				result, err := neo4j.ExecuteQuery(context.Background(), ctx,
					cql,
					map[string]any{
						"bid": id,
						"tid": tag.Id,
					}, neo4j.EagerResultTransformer,
					neo4j.ExecuteQueryWithDatabase("neo4j"))
				fmt.Println(result, err, cql, block.Id, tag.Id, tag.Tag)
			}
		}
		return true
	})
}

/**
  alter 内容

  type HashTag {
      id
      tag
      name
      child
  }



  # Define Directives and index

  id: string @index(term) @lang .
  tag: string @index(term) .
  name: string .
  child: [uid] @count .
*/

func SaveTagsDQL(ctx neo4j.DriverWithContext, tags map[string]*DGraphTag) {
	// 我的 2023 年一个形容词:
	for _, t := range tags {
		cql := fmt.Sprintf("CREATE (n:HashTag {id:$id, name: $name}) RETURN n;")
		result, err := neo4j.ExecuteQuery(context.Background(), ctx,
			cql,
			map[string]any{
				"id":   t.Id,
				"name": t.Tag,
			}, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))
		fmt.Println(result, err, cql, t.Id, t.Tag)
	}
	for _, t := range tags {
		for _, c := range t.Child {
			cql := fmt.Sprintf("MATCH (a:HashTag {id:$fromid}), (b:HashTag {id:$toid}) MERGE (a)-[:CHILD]->(b)")
			result, err := neo4j.ExecuteQuery(context.Background(), ctx,
				cql,
				map[string]any{
					"fromid": t.Id,
					"toid":   c.Id,
				}, neo4j.EagerResultTransformer,
				neo4j.ExecuteQueryWithDatabase("neo4j"))
			fmt.Println(result, err, cql, t.Id, t.Tag, c.Id, c.Tag)
		}
	}
}
