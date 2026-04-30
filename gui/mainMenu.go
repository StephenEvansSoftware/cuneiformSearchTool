package gui

import (
	"fyne.io/fyne/v2"
)

const (
	menuTextOptions = "Options"

	menuTextResults                  = "Results"
	menuItemTextResultsSaveSearch    = "Save search"
	menuItemTextResultsLoadSearch    = "Load search"
	menuItemTextResultsExportResults = "Export results"

	menuTextView           = "View"
	menuItemTextViewSearch = "Search"
	menuItemTextViewCorpus = "Corpus"
)

type mainMenu interface {
	getMainMenu() *fyne.MainMenu
}

type mmStruct struct {
	mainMenu *fyne.MainMenu

	// Menus
	optionsMenu *fyne.Menu

	viewMenu   *fyne.Menu
	viewSearch *fyne.MenuItem
	viewCorpus *fyne.MenuItem

	resultsMenu   *fyne.Menu
	resultsSave   *fyne.MenuItem
	resultsLoad   *fyne.MenuItem
	resultsExport *fyne.MenuItem

	// Reference to parent to pass information back and forth
	gui guiInternal
}

func newMainMenu(gui guiInternal) mainMenu {
	mm := &mmStruct{gui: gui}

	mm.setupOptionsMenu()
	mm.setupViewMenu()
	mm.setupResultsMenu()

	mm.mainMenu = fyne.NewMainMenu(
		mm.optionsMenu,
		mm.viewMenu,
		mm.resultsMenu,
	)

	return mm
}

func (mm *mmStruct) setupOptionsMenu() {
	mm.optionsMenu = fyne.NewMenu(
		menuTextOptions,
	)
}

func (mm *mmStruct) setupViewMenu() {
	// TODO: disable the menu item for views once the view is loaded
	mm.viewSearch = fyne.NewMenuItem(menuItemTextViewSearch, func() { mm.gui.setViewSearch() })
	mm.viewCorpus = fyne.NewMenuItem(menuItemTextViewCorpus, func() { mm.gui.toggleCorpusPane() })

	mm.viewMenu = fyne.NewMenu(
		menuTextView,
		mm.viewSearch,
		mm.viewCorpus,
	)
}

func (mm *mmStruct) setupResultsMenu() {
	mm.resultsSave = fyne.NewMenuItem(menuItemTextResultsSaveSearch, func() {})
	mm.resultsLoad = fyne.NewMenuItem(menuItemTextResultsLoadSearch, func() {})
	mm.resultsExport = fyne.NewMenuItem(menuItemTextResultsExportResults, func() {})

	mm.resultsMenu = fyne.NewMenu(
		menuTextResults,
		mm.resultsSave,
		mm.resultsLoad,
		mm.resultsExport,
	)
}

func (mm *mmStruct) getMainMenu() *fyne.MainMenu {
	return mm.mainMenu
}
