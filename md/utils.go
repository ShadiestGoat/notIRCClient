package md

import (
	"io"
	"math"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func Lines(lines *text.Segments, s []byte) string {
	o := []byte{}

	for i := 0; i < lines.Len(); i++ {
		l := lines.At(i)
		o = append(o, s[l.Start:l.Stop]...)
	}

	return string(o)
}

func logBlock(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	return ast.WalkContinue, nil
}

func logItem(w io.Writer, s []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	return ast.WalkContinue, nil
}

func countDigits(v int) int {
	return int(math.Floor(math.Log10(float64(v))) + 1)
}