package main

import (
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shadiestgoat/log"
)

func init() {
	conf, _ := os.UserConfigDir()
	p := path.Join(conf, "notIRC")
	os.Mkdir(p, 0755)
	p = path.Join(p, "log")

	log.Init(log.NewLoggerFileComplex(p, log.FILE_OVERWRITE, 0))
}

var mainView *MainView
var p *tea.Program

func init() {
	mainView = initMainView()
	p = tea.NewProgram(mainView, 
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
}

func main() {
	_, err := p.Run()
	log.FatalIfErr(err, "running app")
}
