package gui

import (
	"fyne.io/fyne/v2"
)

const (
	menuTextOptions = "Options"

	menuItemTextSaveSearch    = "Save search"
	menuItemTextLoadSearch    = "Load search"
	menuItemTextExportResults = "Export results"
)

type mainMenu interface {
	getMainMenu() *fyne.MainMenu
}

type mmStruct struct {
	mainMenu *fyne.MainMenu

	// Menus
	resultsMenu *fyne.Menu

	// Reference to parent to pass information back and forth
	gui guiInternal
}

func newMainMenu(gui guiInternal) mainMenu {
	mm := &mmStruct{gui: gui}

	mm.resultsMenu = fyne.NewMenu(
		menuTextOptions,
		fyne.NewMenuItem(menuItemTextSaveSearch, func() {}),
		fyne.NewMenuItem(menuItemTextLoadSearch, func() {}),
		fyne.NewMenuItem(menuItemTextExportResults, func() {}),
	)

	mm.mainMenu = fyne.NewMainMenu(mm.resultsMenu)

	return mm
}

func (mm *mmStruct) getMainMenu() *fyne.MainMenu {
	return mm.mainMenu
}
