package md

import (
	"bytes"
	"io"

	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func newParser() goldmark.Markdown {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.NewCJK(),
			extension.Strikethrough,
			emoji.Emoji,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithInlineParsers(
				util.Prioritized(&underlineParser{}, 450),
			),
		),
	)

	r := &Renderer{
		nodeRenderFuncs: map[ast.NodeKind]NodeRendererFunc{},
	}
	
	md.SetRenderer(r)

	return md
}

var mdParser = newParser()

func Parse(inp string) ast.Node {
	return mdParser.Parser().Parse(text.NewReader([]byte(inp)))
}

// baseColor is a hex color w/o the leading #
func RenderText(doc ast.Node, inp string, baseColor string) string {
	r := &Renderer{
		nodeRenderFuncs: map[ast.NodeKind]NodeRendererFunc{},
		bqLevel:         0,
		blockKind:       0,
	}
	r.AddOptions(ColorOpt{
		HexColor: baseColor,
	})

	rw := bytes.NewBuffer(nil)

	r.Render(rw, []byte(inp), doc)

	o, _ := io.ReadAll(rw)
	return string(o)
}

// func Render(inp string) string {
// 	return ""
// }