package main

import (
	"cuneiformSearchTool/gui"
//	"cuneiformSearchTool/util"
)

func main() {
//	// Load paths for all corpus texts
//	if err := util.LoadCorpusTexts(); err != nil {
//		panic(err.Error()) // TODO handle better
//	}
//	gui.TestNewPane()
	gui := gui.NewGUI()
	gui.SetMainWindow()
	gui.ShowGUI()

/*
	// Set up the graphical interface
	gui, err := gui.NewGUI()
	if err != nil {
		panic(err.Error()) // TODO: handle better
	}

	gui.SetMainWindow()
	gui.ShowGUI()
*/
}
