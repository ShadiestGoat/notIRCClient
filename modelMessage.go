package main

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/shadiestgoat/notIRCClient/md"
)

type ModelMessage struct {
	msg *Message

	cachedContent string
	displayName   bool

	ctx *ViewContext
}

func NewModelMessage(msg *Message, ctx *ViewContext) ModelMessage {
	doc := md.Parse(msg.Content)
	cache := md.RenderText(doc, msg.Content, msg.Color())

	shouldDisplayName := msg.Prev == nil || msg.Prev.Author != msg.Author

	return ModelMessage{
		msg:           msg,
		cachedContent: cache,
		ctx:           ctx,
		displayName:   shouldDisplayName,
	}
}

func (m ModelMessage) View() string {
	msg := m.msg
	base := lipgloss.NewStyle().Foreground(lipgloss.Color("#" + msg.Color()))

	contentSize := m.ctx.getContentWidth()

	authorName := ""

	if m.displayName {
		authorName = msg.Author[:len(msg.Author)-6] + ":"

		if m.msg.Prev != nil {
			base = base.MarginTop(1)
		}
	}

	author := styleAuthor.Copy().Width(m.ctx.Width - contentSize - 1).Render(authorName)
	content := lipgloss.NewStyle().Width(contentSize).Render(m.cachedContent)

	return base.Render(
		lipgloss.JoinHorizontal(lipgloss.Top, author, content),
	)
}
