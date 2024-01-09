package store

import (
	"github.com/yuin/goldmark/ast"
	"gopkg.in/yaml.v2"
	"strings"
	"time"
)

type DocumentType int

const (
	D_OTHER DocumentType = iota
	D_MD
	D_PDF
	D_IMG
	D_VIDEO
)

type Document struct {
	Id      string        `json:"id"`
	Path    string        `json:"path"`
	Name    string        `json:"name"`
	Extra   string        `json:"extra"`
	Type    DocumentType  `json:"type"`
	Meta    yaml.MapSlice `json:"meta"`
	Blocks  []*Block      `json:"blocks"`
	ModTime time.Time     `json:"mod_time"`
}

func WalkBlocks(b *Block, walker func(b *Block) bool) {
	if walker(b) {
		for _, c := range b.Children {
			WalkBlocks(c, walker)
		}
	}
}

type Piece struct {
	Id    string
	Doc   Document
	Edges []*Edge
	*Block
}

type Edge struct {
	From *Piece
	To   *Piece
}

type Block struct {
	Level    int                 `json:"level"`
	Id       string              `json:"id"`
	Kind     ast.NodeKind        `json:"kind"`
	HashTags []HashTag           `json:"hash_tags"`
	Inline   map[string][]string `json:"inline"`
	Links    []Link              `json:"links"`
	Content  string              `json:"content"`
	Children []*Block            `json:"children"`
}

type LinkType int

const (
	L_URL LinkType = iota
	L_FILE
	L_IMG
)

type Link struct {
	Type   LinkType
	Target string
	Alias  string
	Show   bool
}

type HashTag []string

// a/b/c -> [a,b,c]
func BuildHashTag(s string) HashTag {
	if s == "" {
		return nil
	}
	return strings.Split(s, "/")
}

func TargetFileName(target string) (name string, path string) {
	if target == "" {
		return "", ""
	}
	if strings.Contains(target, "#") {
		target = strings.Split(target, "#")[0]
	}
	ss := strings.Split(target, "/")
	return ss[0], target
}

type GraphNode struct {
	Id       string
	Doc      *Document
	Block    *Block
	OutEdges []GraphEdge
}

type GraphEdge struct {
	From  string
	To    string
	Alias string
}
