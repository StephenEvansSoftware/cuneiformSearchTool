package gui

import (
	"fmt"

	"fyne.io/fyne/v2"

	"cuneiformSearchTool/util"
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
	// TODO: quick hacky way to print out the corpus meta info, fix
	optionsPrintMetaInfo := fyne.NewMenuItem("Print Corpus Meta Info", func() {
		corpusOrder := []string{}
		corpusMetaInfo := map[string]*util.CorpusMetaInfo{}
		for _, ci := range util.ListCorpusAndTexts() {
			corpusOrder = append(corpusOrder, ci.CorpusName)
			cmi, err := util.GetCorpusMetaInfo(ci.CorpusName)
			if err != nil {
				panic(err.Error()) // TODO: handle better
			}
			corpusMetaInfo[ci.CorpusName] = cmi
		}

		for _, co := range corpusOrder {
			fmt.Printf(",%s", co)
		}
		fmt.Printf("\nNumberOfTexts")
		for _, co := range corpusOrder {
			fmt.Printf(",%d", corpusMetaInfo[co].NumberOfTexts)
		}
		fmt.Printf("\nTotalFragsInAllTexts")
		for _, co := range corpusOrder {
			fmt.Printf(",%d", corpusMetaInfo[co].TotalFragsInAllTexts)
		}
		fmt.Printf("\nTotalFragsInAllTextsWithoutX")
		for _, co := range corpusOrder {
			fmt.Printf(",%d", corpusMetaInfo[co].TotalFragsInAllTextsWithoutX)
		}
		fmt.Printf("\nTotalLinesInAllTexts")
		for _, co := range corpusOrder {
			fmt.Printf(",%d", corpusMetaInfo[co].TotalLinesInAllTexts)
		}
		fmt.Printf("\nAverageLineLength")
		for _, co := range corpusOrder {
			fmt.Printf(",%.2f", corpusMetaInfo[co].AverageLineLength)
		}
		fmt.Printf("\nAverageLineLengthWithoutBlanks")
		for _, co := range corpusOrder {
			fmt.Printf(",%.2f", corpusMetaInfo[co].AverageLineLengthWithoutBlanks)
		}
		fmt.Println()
	})

	mm.optionsMenu = fyne.NewMenu(
		menuTextOptions,
		optionsPrintMetaInfo,
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
