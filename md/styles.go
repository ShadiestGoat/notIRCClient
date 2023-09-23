package md

import (
	"io"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/shadiestgoat/colorutils"
	"github.com/shadiestgoat/log"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
)

const (
	COLOR_KEY = "COLOR"
	BORDER_COLOR_KEY = "BORDER_COLOR"
)

type ColorOpt struct {
	// A hex color without the leading #
	HexColor string 
}

func (c ColorOpt) SetConfig(cfg *renderer.Config) {
	if len(c.HexColor) != 6 {
		log.Fatal("Non 6 color hex!")
	}
	rR := c.HexColor[:2]
	rG := c.HexColor[2:4]
	rB := c.HexColor[4:]
	
	r, err := strconv.ParseInt(rR, 16, 0)
	log.FatalIfErr(err, "parsing r")
	g, err := strconv.ParseInt(rG, 16, 0)
	log.FatalIfErr(err, "parsing g")
	b, err := strconv.ParseInt(rB, 16, 0)
	log.FatalIfErr(err, "parsing b")

	h, _, _ := colorutils.RGBToHSL(uint8(r), uint8(g), uint8(b))
	
	cfg.Options[COLOR_KEY] = lipgloss.Color("#" + c.HexColor)
	cfg.Options[BORDER_COLOR_KEY] = lipgloss.Color("#" + colorutils.Hexadecimal(
		colorutils.HSLToRGB(h, 0.3, 0.6),
	))
}

var (
	styleBold   = lipgloss.NewStyle().Bold(true)
	styleStrike = lipgloss.NewStyle().Strikethrough(true)
	styleBlockQuote = lipgloss.NewStyle().Border(lipgloss.Border{
		Left:        "▌",
	}, false, false, false, true).BorderForeground(lipgloss.Color("#4E5058"))
)

func renderStyle(s lipgloss.Style, w io.Writer, n ast.Node, source []byte) {
	w.Write([]byte(s.Render(
		string(n.Text(source)),
	)))
	lipgloss.NormalBorder()
}

func factoryRenderStyle(style lipgloss.Style) NodeRendererFunc {
	return func(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			renderStyle(style, w, n, s)
		}

		return ast.WalkContinue, nil
	}
}
