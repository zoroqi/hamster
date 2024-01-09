package goldmark_obsidian

import (
	"bytes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"unicode"
)

// Kind is the kind of the wikilink AST node.
var BlockIDKind = ast.NewNodeKind("ObsidianBlockId")

type BlockIdNode struct {
	ast.BaseInline

	// BlockId of the Obsidian block.
	// like ^abc
	BlockId []byte
}

var _ ast.Node = (*BlockIdNode)(nil)

// Kind reports the kind of this node.
func (n *BlockIdNode) Kind() ast.NodeKind {
	return BlockIDKind
}

// Dump dumps the BlockIdNode to stdout.
func (n *BlockIdNode) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, map[string]string{
		"BlockId": string(n.BlockId),
	}, nil)
}

// [Internal links - Obsidian Help](https://help.obsidian.md/Linking+notes+and+files/Internal+links#Link+to+a+block+in+a+note)
type BlockParser struct {
}

var _ parser.InlineParser = (*BlockParser)(nil)

var (
	_blockid_open = []byte("^")
)

// Trigger returns characters that trigger this parser.
func (p *BlockParser) Trigger() []byte {
	return _blockid_open
}

func (p *BlockParser) Parse(node ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()

	// The charsets must be "a-z", "A-Z", "0-9" and "-"
	if bytes.ContainsFunc(line[1:], func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-')
	}) {
		return nil
	}

	// prev rune is must space or \n
	prev := block.PrecendingCharacter()
	if !(prev == ' ' || prev == '\n' || prev == '\r') {
		return nil
	}

	n := &BlockIdNode{BlockId: line[1:]}
	block.Advance(len(line))
	return n
}

type BlockExtender struct {
}

var _ goldmark.Extender = (*BlockExtender)(nil)

// Extend extends the provided goldmark Markdown object with support for
// hashtags.
func (e *BlockExtender) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(&BlockParser{}, 999),
		),
	)
}
