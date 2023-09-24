package md

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Underline = basically stolen code from emphasis, but placed at higher priority and only with _

type Underline struct {
	ast.BaseInline

	Level int
}

func (n *Underline) Dump(source []byte, level int) {
	m := map[string]string{
		"Level": fmt.Sprint(n.Level),
	}
	ast.DumpHelper(n, source, level, m, nil)
}

var KindUnderline = ast.NewNodeKind("Underline")

// Kind implements Node.Kind.
func (n *Underline) Kind() ast.NodeKind {
	return KindUnderline
}

type underlineDelimiterProcessor struct {
}

func (p *underlineDelimiterProcessor) IsDelimiter(b byte) bool {
	return b == '_'
}

func (p *underlineDelimiterProcessor) CanOpenCloser(opener, closer *parser.Delimiter) bool {
	return opener.Char == closer.Char
}

func (p *underlineDelimiterProcessor) OnMatch(consumes int) ast.Node {
	return &Underline{
		BaseInline: ast.BaseInline{},
		Level: consumes,
	}
}

type underlineParser struct {
}

func (s *underlineParser) Trigger() []byte {
	return []byte{'_'}
}

func (s *underlineParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	before := block.PrecendingCharacter()
	line, segment := block.PeekLine()
	node := parser.ScanDelimiter(line, before, 1, &underlineDelimiterProcessor{})
	if node == nil {
		return nil
	}
	node.Segment = segment.WithStop(segment.Start + node.OriginalLength)
	block.Advance(node.OriginalLength)
	pc.PushDelimiter(node)
	return node
}

