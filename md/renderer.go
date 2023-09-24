package md

import (
	"bytes"
	"io"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"

	east "github.com/yuin/goldmark-emoji/ast"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type NodeRendererFunc func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error)

type Renderer struct {
	cfg *renderer.Config
	nodeRenderFuncs map[ast.NodeKind]NodeRendererFunc

	baseS lipgloss.Style
	s lipgloss.Style
	bqStyle lipgloss.Style

	bqLevel int

	blockKind ast.NodeKind
}

type Debug struct {
	Text     string
	Entering bool
	Kind     string
}

type StyleWrapper func (bool) lipgloss.Style

func (r *Renderer) UpdateStyle(e bool, f StyleWrapper) {
	r.s = f(e)
}

func (r *Renderer) updateStyle(f StyleWrapper, badParents...ast.NodeKind) NodeRendererFunc {
	return func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		for _, b := range badParents {
			if b == r.blockKind {
				return ast.WalkContinue, nil
			}
		}

		r.s = f(entering)

		return ast.WalkContinue, nil
	}
}

func (r *Renderer) renderText(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.renderBaseText(w, n.Text(s))
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderBaseText(w io.Writer, s []byte) {
	w.Write(
		[]byte(
			r.s.Render(string(s)),
		),
	)
}

func (r *Renderer) renderUnstyledText(w io.Writer, s []byte) {
	w.Write(
		[]byte(
			r.baseS.Render(string(s)),
		),
	)
}

func (r *Renderer) renderRawText(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.renderUnstyledText(w, n.Text(s))
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderBlock(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	r.renderUnstyledText(w, []byte(strings.TrimSuffix(Lines(n.Lines(), s), "\n")))

	return ast.WalkSkipChildren, nil
}

func (r *Renderer) RegisterRenderers() {
	r.baseS = lipgloss.NewStyle().StrikethroughSpaces(true).UnderlineSpaces(true)
	r.bqStyle = styleBlockQuote.Copy()
	
	if r.cfg != nil {
		if opt, ok := r.cfg.Options[COLOR_KEY]; ok {
			r.baseS = r.baseS.Foreground(opt.(lipgloss.Color))
			r.bqStyle = r.bqStyle.BorderForeground(r.cfg.Options[BORDER_COLOR_KEY].(lipgloss.Color))
		}
	}

	r.s = r.baseS.Copy()

	var (
		renderEmpty = func(w io.Writer, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
			return ast.WalkContinue, nil
		}

		// renderEmptyContent = func(w io.Writer, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		// 	return ast.WalkSkipChildren, nil
		// }

		// renderBold = factoryRenderStyle(styleBold)
	)

	r.nodeRenderFuncs[ast.KindDocument] = renderEmpty
	r.nodeRenderFuncs[ast.KindParagraph] = renderEmpty
	r.nodeRenderFuncs[ast.KindList] = renderEmpty
	
	r.nodeRenderFuncs[ast.KindHeading] = func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			h := n.(*ast.Heading)

			b := []byte{}
			for i := 0; i < h.Level - 1; i++ {
				b = append(b, ' ')
			}
			
			w.Write(b)
		}

		r.s = r.s.Bold(entering).Underline(entering)

		return ast.WalkContinue, nil
	}

	r.nodeRenderFuncs[extast.KindStrikethrough] = r.updateStyle(r.s.Strikethrough)
	
	r.nodeRenderFuncs[ast.KindFencedCodeBlock] = func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		block := n.(*ast.FencedCodeBlock)

		rawText := "```"

		if block.Info != nil {
			rawText += string(block.Info.Text(s))
		}

		rawText += "\n" + Lines(block.BaseBlock.Lines(), s) + "```"
		
		r.renderBaseText(w, []byte(rawText))
		
		return ast.WalkContinue, nil
	}
	r.nodeRenderFuncs[ast.KindCodeBlock] = func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		r.renderBaseText(w, []byte(
			"```\n" + Lines(n.Lines(), s) + "```",
		))
		
		return ast.WalkContinue, nil
	}

	r.nodeRenderFuncs[ast.KindHTMLBlock] = r.renderBlock

	r.nodeRenderFuncs[ast.KindRawHTML] = func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		segs := n.(*ast.RawHTML).Segments

		if segs.Len() == 0 {
			return ast.WalkSkipChildren, nil
		}

		seg := segs.At(0)

		r.renderBaseText(w, seg.Value(s))

		return ast.WalkSkipChildren, nil
	}

	r.nodeRenderFuncs[ast.KindCodeSpan] = func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			return ast.WalkSkipChildren, nil
		}
		
		r.renderUnstyledText(w, n.FirstChild().Text(s))
		
		return ast.WalkContinue, nil
	}

	r.nodeRenderFuncs[ast.KindBlockquote] = func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			r.bqLevel++
		} else {
			r.bqLevel--
		}

		return ast.WalkContinue, nil
	}

	r.nodeRenderFuncs[ast.KindListItem] = func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			li := n.(*ast.ListItem)
			l := n.Parent().(*ast.List)
			b := []byte{}
			
			suffix := ""
			suffixOffset := 0
			
			if l.IsOrdered() {
				maxChild := l.ChildCount() + l.Start

				i := 0
				prev := li.PreviousSibling()

				for prev != nil {
					prev = prev.PreviousSibling()
					i++
				}

				ind := l.Start + i
				suffix = strconv.Itoa(l.Start + i)
				suffixOffset = countDigits(maxChild) - countDigits(ind)
			}

			for i := 0; i < li.Offset - 2 - len(suffix) + suffixOffset; i++ {
				b = append(b, ' ')
			}

			b = append(b, []byte(suffix)...)
			
			r.renderUnstyledText(w, append(b, []byte{l.Marker, ' '}...))
		}

		return ast.WalkContinue, nil
	}
	

	r.nodeRenderFuncs[ast.KindTextBlock] = renderEmpty

	r.nodeRenderFuncs[ast.KindThematicBreak] = renderEmpty

	r.nodeRenderFuncs[ast.KindAutoLink] = r.renderRawText

	r.nodeRenderFuncs[ast.KindLink] = func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		l := n.(*ast.Link)
		special := l.ChildCount() == 1 && l.FirstChild().Kind() == ast.KindText && string(l.FirstChild().Text(s)) == string(l.Destination)

		if entering {
			if special {
				r.renderUnstyledText(w, l.Destination)
				return ast.WalkSkipChildren, nil
			}

			r.renderUnstyledText(w, []byte{'['})
		} else if !special {
			b := append([]byte{']', '('}, l.Destination...)
			b = append(b, ')')

			r.renderUnstyledText(w, b)
		}

		return ast.WalkContinue, nil
	}

	r.nodeRenderFuncs[ast.KindImage] = func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		l := n.(*ast.Image)

		if entering {
			r.renderUnstyledText(w, []byte{'!', '['})
		} else {
			b := append([]byte{']', '('}, l.Destination...)
			b = append(b, ')')
			
			r.renderUnstyledText(w, b)
		}

		return ast.WalkContinue, nil
	}

	r.nodeRenderFuncs[ast.KindEmphasis] = func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		e := n.(*ast.Emphasis)

		var f NodeRendererFunc

		switch e.Level {
		case 1:
			f = r.updateStyle(r.s.Italic)
		case 2:
			f = r.updateStyle(r.s.Bold, ast.KindHeading)
		}

		if f != nil {
			f(w, s, n, entering)
		}

		return ast.WalkContinue, nil
	}

	r.nodeRenderFuncs[KindUnderline] = func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		e := n.(*Underline)

		var f NodeRendererFunc

		switch e.Level {
		case 1:
			f = r.updateStyle(r.s.Italic)
		case 2:
			f = r.updateStyle(r.s.Underline, ast.KindHeading)
		}

		if f != nil {
			f(w, s, n, entering)
		}

		return ast.WalkContinue, nil
	}

	r.nodeRenderFuncs[ast.KindText] = func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		t := n.(*ast.Text)

		if entering {
			r.renderText(w, s, n, entering)
		} else if t.HardLineBreak() || t.SoftLineBreak() {
			w.Write([]byte{'\n'})
		}
		
		return ast.WalkContinue, nil
	}
	r.nodeRenderFuncs[ast.KindString] = r.renderText

	r.nodeRenderFuncs[extast.KindTaskCheckBox] = func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		
		li := n.(*extast.TaskCheckBox)

		checkStatus := " "
		
		if li.IsChecked {
			checkStatus = "✅"
		}

		r.renderUnstyledText(w, []byte("[" + checkStatus + "] "))

		return ast.WalkContinue, nil
	}

	r.nodeRenderFuncs[east.KindEmoji] = func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			emote := n.(*east.Emoji)
			w.Write([]byte(string(emote.Value.Unicode)))
		}

		return ast.WalkContinue, nil 
	}
}

func (r *Renderer) Render(finalW io.Writer, source []byte, n ast.Node) error {
	var bqW = []io.Writer{
		finalW,
	}

	r.RegisterRenderers()

	lastBQLvl := 0

	err := ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		var err error
		f := r.nodeRenderFuncs[n.Kind()]

		blockType := n.Type() == ast.TypeBlock
		rootBlock := false

		if blockType {
			r.blockKind = n.Kind()

			firstChild := n.PreviousSibling() == nil
			rootBlock = n.Parent().Type() == ast.TypeDocument
			
			if entering && n.HasBlankPreviousLines() && !firstChild && rootBlock {
				bqW[r.bqLevel].Write([]byte("\n"))
			}
		}

		if f != nil {
			s, err = f(bqW[r.bqLevel], source, n, entering)
		}

		isLastChild := n.NextSibling() == nil
		
		if r.bqLevel > lastBQLvl {
			bqW = append(bqW, bytes.NewBuffer(nil))
		}

		if r.bqLevel < lastBQLvl {
			b, _ := io.ReadAll(bqW[lastBQLvl].(io.ReadWriter))

			bqW[r.bqLevel].Write(
				[]byte(r.bqStyle.Render(string(b))),
			)
		}

		if blockType && (rootBlock || !isLastChild) && !(rootBlock && isLastChild) {
			if !entering {
				bqW[r.bqLevel].Write([]byte("\n"))
			}
		}

		lastBQLvl = r.bqLevel 

		return s, err
	})

	return err
}

func (r *Renderer) AddOptions(opts ...renderer.Option) {
	if r.cfg == nil {
		r.cfg = &renderer.Config{
			Options:       map[renderer.OptionName]interface{}{},
			NodeRenderers: []util.PrioritizedValue{},
		}
	}

	for _, o := range opts {
		o.SetConfig(r.cfg)
	}
}
