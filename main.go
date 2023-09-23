package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shadiestgoat/log"
)

func init() {
	log.Init(log.NewLoggerFile("log"))
}

var mainView = initMainView()
var p = tea.NewProgram(mainView, 
	tea.WithAltScreen(),
	tea.WithMouseCellMotion(),
)

func main() {
	_, err := p.Run()
	log.FatalIfErr(err, "running app")
}
