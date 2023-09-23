package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
	"github.com/shadiestgoat/log"
)

type MainView struct {
	c *websocket.Conn

	msgCache string
	msgs []ModelMessage

	// Author ID (inc hex)
	authorCache map[string]bool
	ctx *ViewContext

	help help.Model
	viewport viewport.Model
	textarea textarea.Model
}

func initMainView() *MainView {
	conn, _, err := websocket.DefaultDialer.Dial(urlBase("ws") + "/ws", nil)
	log.FatalIfErr(err, "dialing ws")

	messages := GetMessages()
	
	authors := map[string]bool{}
	msgs := []ModelMessage{}

	ctx := &ViewContext{}

	for i, m := range messages {
		authors[m.Author] = true
		
		if i != 0 {
			m.Prev = messages[i - 1]
		}

		msgs = append(msgs, NewModelMessage(m, ctx))
	}

	for a := range authors {
		s := len(a) - 6
		if s > ctx.AuthorSize {
			ctx.AuthorSize = s
		}
	}

	return &MainView{
		c:           conn,
		msgs:        msgs,
		authorCache: authors,
		ctx:  ctx,
	}
}

func (m *MainView) renderMessages() {
	views := []string{}

	for _, msg := range m.msgs {
		v := msg.View()
		views = append(views, v)
	}

	m.msgCache = lipgloss.JoinVertical(0, views...)
}

func (m *MainView) Init() tea.Cmd {
	go func() {
		var (
			msg *Message
		)

		for {
			err := m.c.ReadJSON(&msg)
			if log.ErrorIfErr(err, "reading WS json") {
				p.Quit()
				return
			}

			m.Update(msg)
		}
	}()

	return nil
}

func (m *MainView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd

		oldCtx = *m.ctx

		dontUpdateTextArea bool
		shouldScrollToBot bool
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Debug(msg.String())

		switch {
		case key.Matches(msg, appKeyMap.Quit):
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, appKeyMap.SendMsg):
			content := strings.TrimSpace(m.textarea.Value())
			if content == "" {
				return m, nil
			}

			SendMessage(m.textarea.Value())
			m.textarea.Reset()
			dontUpdateTextArea = true
		case key.Matches(msg, appKeyMap.Help):
			m.help.ShowAll = !m.help.ShowAll
			dontUpdateTextArea = true

			hCount := strings.Count(m.help.View(appKeyMap), "\n") + 1

			oldH := m.viewport.Height
			newH := m.ctx.Height - TEXT_INPUT_HEIGHT - hCount
			
			m.viewport.Height = newH
			m.viewport.SetYOffset(m.viewport.YOffset + (oldH - newH))
		}
	case *Message:
		if !m.authorCache[msg.Author] {
			m.authorCache[msg.Author] = true
			s := len(msg.Author) - 6

			if s > m.ctx.AuthorSize {
				m.ctx.AuthorSize = s
			}
		}

		if len(m.msgs) != 0 {
			msg.Prev = m.msgs[len(m.msgs)-1].msg
		}

		msgMod := NewModelMessage(msg, m.ctx)

		m.msgs = append(m.msgs, msgMod)

		if oldCtx.AuthorSize == 0 {

		} else if m.ctx.AuthorSize == oldCtx.AuthorSize {
			m.msgCache = lipgloss.JoinVertical(0, m.msgCache, msgMod.View())
		} else {
			m.renderMessages()
		}

		shouldScrollToBot = msg.Author == AUTHOR_SEND_INFO
	case tea.WindowSizeMsg:
		m.ctx.Width = msg.Width
		m.ctx.Height = msg.Height

		if oldCtx.Width != msg.Width {
			m.renderMessages()
		}

		helpHeight := 1

		if oldCtx.Width != 0 {
			helpHeight = strings.Count(m.help.View(appKeyMap), "\n") + 1
		}

		vpHeight := msg.Height - TEXT_INPUT_HEIGHT - helpHeight

		if oldCtx.Width == 0 {
			m.viewport = viewport.New(msg.Width, vpHeight)
			m.viewport.KeyMap = vpKeys

			m.textarea = textarea.New()
			m.textarea.KeyMap = textAreaKeys
			m.textarea.ShowLineNumbers = false
			m.textarea.Placeholder = "Message..."
			m.textarea.SetHeight(TEXT_INPUT_HEIGHT)
			m.textarea.SetWidth(msg.Width)
			m.textarea.Focus()
			
			m.help = help.New()
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = vpHeight
			m.textarea.SetWidth(msg.Width)
		}
		
		m.help.Width = msg.Width
	}

	m.viewport.SetContent(m.msgCache)

	if shouldScrollToBot {
		m.viewport.GotoBottom()
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	if _, ok := msg.(tea.MouseMsg); !ok && !dontUpdateTextArea {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m MainView) View() string {
	if m.ctx.Width == 0 {
		return "Loading..."
	}
	log.Debug("View called")

	return lipgloss.JoinVertical(0, 
		m.viewport.View(),
		m.help.View(appKeyMap),
		m.textarea.View(),
	)
}
