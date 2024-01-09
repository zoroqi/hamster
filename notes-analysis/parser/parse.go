package parser

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	obsidian "github.com/zoroqi/hamster/notes-analysis/parser/goldmark-obsidian"
	"github.com/zoroqi/hamster/notes-analysis/store"
	"go.abhg.dev/goldmark/hashtag"
	"go.abhg.dev/goldmark/wikilink"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileInfo struct {
	os.FileInfo
	Path string
}

func ParseAllfile(root string) ([]*store.Document, error) {
	files := []FileInfo{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// filter dot file
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}
		if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
			files = append(files, FileInfo{info, path})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	docs := make([]*store.Document, 0, len(files))
	for _, file := range files {
		d := createDoc(file)
		docs = append(docs, d)
	}

	for _, doc := range docs {
		if doc.Type == store.D_MD {
			f, err := os.Open(doc.Path)
			if err != nil {
				fmt.Println(doc.Path, err)
				return nil, err
			}
			meta, blocks, err := ParseMarkdown(f, doc)
			if err != nil {
				return nil, err
			}
			doc.Meta = meta
			doc.Blocks = blocks
		}
	}

	return docs, nil

}

func createEmptyDoc(file string) *store.Document {
	fileid := func(path string) string {
		m1 := md5.Sum([]byte(path))
		return fmt.Sprintf("%x", m1)
	}
	extra := filepath.Ext(file)
	name := filepath.Base(file)
	name = name[:len(name)-len(extra)]
	if extra == ".png" || extra == ".jpg" || extra == ".jpeg" || extra == ".gif" {
		return &store.Document{
			Id:      fileid(file),
			Path:    file,
			Type:    store.D_IMG,
			Name:    name,
			Extra:   extra,
			ModTime: time.Now(),
		}
	} else if extra == ".pdf" {
		return &store.Document{
			Id:      fileid(file),
			Path:    file,
			Type:    store.D_PDF,
			Name:    name,
			Extra:   extra,
			ModTime: time.Now(),
		}
	} else if extra == ".mp4" || extra == ".avi" || extra == ".mkv" {
		return &store.Document{
			Id:      fileid(file),
			Path:    file,
			Type:    store.D_VIDEO,
			Name:    name,
			Extra:   extra,
			ModTime: time.Now(),
		}
	} else if extra == ".md" {
		return &store.Document{
			Id:      fileid(file),
			Path:    file,
			Type:    store.D_MD,
			Name:    name,
			Extra:   extra,
			ModTime: time.Now(),
		}
	} else {
		return &store.Document{
			Id:      fileid(file),
			Path:    file,
			Type:    store.D_OTHER,
			Name:    name,
			Extra:   extra,
			ModTime: time.Now(),
		}
	}
}

func createDoc(info FileInfo) *store.Document {
	doc := createEmptyDoc(info.Path)
	doc.ModTime = info.ModTime()
	return doc
}

func ParseMarkdown(read io.Reader, doc *store.Document) (yaml.MapSlice, []*store.Block, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			//extension.TaskList,
			meta.Meta,
			&wikilink.Extender{},
			&hashtag.Extender{},
			&obsidian.BlockExtender{},
			&obsidian.InlineFieldsExtender{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	path := doc.Path

	source, err := io.ReadAll(read)
	if err != nil {
		return nil, nil, err
	}
	context := parser.NewContext()
	reader := text.NewReader(source)
	md5h := md5.New()
	id := func(n ast.Node) string {
		k := base64.StdEncoding.EncodeToString(md5h.Sum(n.Text(source)))
		m1 := md5.Sum([]byte(path + "_" + n.Kind().String() + k))
		return fmt.Sprintf("%x", m1)[0:8]
	}

	blockKind := map[ast.NodeKind]bool{
		ast.KindHeading:         true,
		ast.KindParagraph:       true,
		ast.KindBlockquote:      true,
		ast.KindList:            true,
		ast.KindListItem:        true,
		east.KindTable:          true,
		ast.KindFencedCodeBlock: true,
	}
	valueKind := map[ast.NodeKind]bool{
		obsidian.InlineFieldsKind: true,
		obsidian.BlockIDKind:      true,
		ast.KindLink:              true,
		wikilink.Kind:             true,
		ast.KindTextBlock:         true,
		ast.KindImage:             true,
		hashtag.Kind:              true,
	}

	root := md.Parser().Parse(reader, parser.WithContext(context))
	var blockDfs func(n ast.Node, dep int) []*store.Block
	var valueDfs func(n ast.Node, dep int, block *store.Block)
	valueDfs = func(n ast.Node, dep int, block *store.Block) {
		switch v := n.(type) {
		case *obsidian.InlineFieldsNode:
			k := string(v.Key)
			value := string(v.Value)
			block.Inline[k] = append(block.Inline[k], value)
		case *ast.Link:
			target := string(v.Destination)
			title := string(v.Title)
			block.Links = append(block.Links,
				store.Link{Type: store.L_URL, Target: target, Alias: title, Show: false})
		case *ast.Image:
			target := string(v.Destination)
			title := string(v.Title)
			block.Links = append(block.Links,
				store.Link{Type: store.L_IMG, Target: target, Alias: title, Show: true})
		case *wikilink.Node:
			target := string(v.Target)
			if target == "" {
				target = doc.Name
			}
			title := string(v.Fragment)
			block.Links = append(block.Links,
				store.Link{Type: store.L_FILE, Target: target,
					Alias: title, Show: v.Embed})
		case *obsidian.BlockIdNode:
			block.Id = string(v.BlockId) + "_" + block.Id[0:8]
		case *hashtag.Node:
			target := string(v.Tag)
			block.HashTags = append(block.HashTags, store.BuildHashTag(target))
		}
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			valueDfs(c, dep+1, block)
		}
	}
	blockDfs = func(n ast.Node, dep int) []*store.Block {
		b := &store.Block{
			Level:   dep,
			Id:      id(n),
			Content: string(n.Text(source)),
			Inline:  make(map[string][]string),
			Kind:    n.Kind(),
		}
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			switch {
			case blockKind[c.Kind()]:
				cbs := blockDfs(c, dep+1)
				b.Children = append(b.Children, cbs...)
			case valueKind[c.Kind()]:
				valueDfs(c, dep+1, b)
			default:
			}

		}
		return []*store.Block{b}
	}
	blocks := blockDfs(root, 0)

	var relatedBlockId func(block *store.Block)

	relatedBlockId = func(block *store.Block) {
		for i := 1; i < len(block.Children); i++ {
			b := block.Children[i]
			if b.Kind == ast.KindParagraph && strings.TrimSpace(b.Content) == "" {
				//fmt.Println("--------")
				block.Children[i-1].Id = b.Id
				block.Children = append(block.Children[:i], block.Children[i+1:]...)
				i--
			} else {
				relatedBlockId(block.Children[i])
			}
		}
	}

	//var moveHeadChange func(blocks *store.Block)
	//moveHeadChange = func(b *store.Block) {
	//	blocks := b.Child
	//	for i := 0; i < len(blocks); i++ {
	//		if blocks[i].Kind == ast.KindHeading {
	//			for j := i + 1; j < len(blocks); j++ {
	//				if blocks[j].Kind != ast.KindHeading {
	//					break
	//				}
	//				blocks[i].Child = append(blocks[i].Child, blocks[j])
	//			}
	//		}
	//	}
	//}

	var dfsBlock func(*store.Block, int)
	dfsBlock = func(block *store.Block, i int) {
		//fmt.Printf("%s%d,%s,%s,%s,%d,%d,%v,%s\n", strings.Repeat("  ", i),
		//	block.Level, block.Kind.String(), block.Id,
		//	block.Content[0:min(12, len(block.Content))], len(block.Links), len(block.Child), block.Inline, block.HashTags)
		for _, v := range block.Children {
			dfsBlock(v, i+1)
		}
	}
	for _, v := range blocks {
		relatedBlockId(v)
	}
	for _, v := range blocks {
		dfsBlock(v, 0)
	}

	metaData := meta.GetItems(context)
	//doc := &store.Document{
	//	Meta: metaData,
	//}
	return metaData, blocks, nil

}

func printNode(dep int, n ast.Node, bs string, v string) {
	fmt.Printf("%s%s:%s:%s\n", strings.Repeat("  ", dep), v, n.Kind(), bs)
}
