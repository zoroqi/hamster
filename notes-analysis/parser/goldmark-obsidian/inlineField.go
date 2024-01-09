package goldmark_obsidian

import (
	"bytes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// [Adding Metadata - Dataview](https://blacksmithgu.github.io/obsidian-dataview/annotation/add-metadata/#inline-fields)
// For those wanting a more natural-looking annotation,
// Dataview supports "inline" fields via a Key:: Value syntax that you can use everywhere in your file.
// This allows you to write your queryable data right where you need it - for example in the middle of a sentence.
type InlineFieldsNode struct {
	ast.BaseInline
	Key   []byte
	Value []byte
}

var InlineFieldsKind = ast.NewNodeKind("InlineFields")

var _ ast.Node = (*InlineFieldsNode)(nil)

// Kind reports the kind of this node.
func (n *InlineFieldsNode) Kind() ast.NodeKind {
	return InlineFieldsKind
}

// Dump dumps the BlockIdNode to stdout.
func (n *InlineFieldsNode) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, map[string]string{
		"key":   string(n.Key),
		"value": string(n.Value),
	}, nil)
}

type InlineFieldsParser struct {
}

var _ parser.InlineParser = (*InlineFieldsParser)(nil)

var (
	_inline_fields_open  = []byte("[")
	_inline_fields_close = []byte("]")
	_inline_field_middle = []byte("::")

	_wiki_open = []byte("[[")
)

// Trigger returns characters that trigger this parser.
func (p *InlineFieldsParser) Trigger() []byte {
	return []byte("[")
}

func (p *InlineFieldsParser) Parse(node ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, seg := block.PeekLine()
	if bytes.HasPrefix(line, _wiki_open) {
		return nil
	}

	stop := bytes.Index(line, _inline_fields_close)
	if stop < 0 {
		return nil
	}

	middle := bytes.Index(line, _inline_field_middle)
	// [::abc]
	if middle-1 <= 0 {
		return nil
	}
	if stop-middle-1 <= 0 {
		return nil
	}

	key := line[1:middle]
	value := line[middle+2 : stop]

	n := &InlineFieldsNode{Key: key, Value: value}
	seg = text.NewSegment(seg.Start, seg.Start+stop+len(_inline_fields_close))
	n.AppendChild(n, ast.NewTextSegment(seg))
	block.Advance(seg.Len())
	return n
}

type InlineFieldsExtender struct {
}

var _ goldmark.Extender = (*InlineFieldsExtender)(nil)

// Extend extends the provided goldmark Markdown object with support for
// hashtags.
func (e *InlineFieldsExtender) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			//
			util.Prioritized(&InlineFieldsParser{}, 198),
		),
	)
}
