package main

import (
)


type ViewContext struct {
	Width int
	Height int

	// The biggest author NAME size
	AuthorSize int
}

func (c ViewContext) getContentWidth() int {
	return c.Width - c.AuthorSize - 2
}
