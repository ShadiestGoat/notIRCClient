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

	displayName := msg.Prev == nil || !msg.Prev.lastAuthorUnderX(msg.Author, 7)

	return ModelMessage{
		msg:           msg,
		cachedContent: cache,
		ctx:           ctx,
		displayName:   displayName,
	}
}

func (m Message) lastAuthorUnderX(searchName string, needed int) bool {
	if m.Author == searchName {
		return true
	}

	if needed <= 0 {
		return false
	}

	if m.Prev == nil {
		return false
	}

	return m.Prev.lastAuthorUnderX(searchName, needed-1)
}

func (m ModelMessage) shouldDisplayName() bool {
	return !m.msg.lastAuthorUnderX(m.msg.Author, 7)
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
