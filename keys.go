package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
)

var vpKeys = viewport.KeyMap{
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	HalfPageUp: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "½ page up"),
	),
	HalfPageDown: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "½ page down"),
	),
	Up: key.NewBinding(
		key.WithKeys("ctrl+k"),
		key.WithHelp("Ctrl+K/Scroll", "Scroll up"),
	),
	Down: key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("Ctrl+J/Scroll", "Scroll down"),
	),
}

var textAreaKeys = textarea.DefaultKeyMap

func init() {
	textAreaKeys.DeleteAfterCursor = key.NewBinding()
	textAreaKeys.DeleteBeforeCursor = key.NewBinding()
	textAreaKeys.DeleteCharacterBackward = key.NewBinding(key.WithKeys("backspace"))
	textAreaKeys.DeleteCharacterForward = key.NewBinding(key.WithKeys("delete"))
	textAreaKeys.InsertNewline = key.NewBinding(
		key.WithKeys("alt+enter", "alt+'"),
		key.WithHelp("Alt+'/Alt+Enter", "Pular linha"),
	)

	textAreaKeys.Paste.SetHelp("Ctrl+V", "Colar")
}

type fullKeyMap struct {
	Help,
	VPUp,
	VPDown,
	NewLine,
	SendMsg,
	Copy,
	Paste,
	Quit key.Binding
}

func (k fullKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k fullKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.SendMsg, k.Copy},
		{k.Quit, k.NewLine, k.Paste},
	}
}

var appKeyMap fullKeyMap


func init() {
	appKeyMap = fullKeyMap{
		Help: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("Ctrl+H", "Abrir menu de atalhos"),
		),
		VPUp:       vpKeys.Up,
		VPDown:     vpKeys.Down,
		NewLine:    textAreaKeys.InsertNewline,
		SendMsg: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("Enter", "Enviar mensagem"),
		),
		Copy:       key.NewBinding(
			key.WithKeys("ctrl+shift+c"),
			key.WithHelp("Ctrl+Shift+C", "Copiar"),
		),
		Paste:      textAreaKeys.Paste,
		Quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("ESC", "Sair"),
		),
	}
}