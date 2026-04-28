package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"cuneiformSearchTool/util"
)

var (
	// Version is injected at compile time
	VERSION_MAJOR, VERSION_MINOR, VERSION_PATCH string
	VERSION                                     string
)

const (
	APP_ID            = "cunfeiform-search-tool"
	MAIN_WINDOW_TITLE = "Cuneiform Search Tool"

	// TODO: hack to force list being fully visible in accordion item inside vbox
	SIZE_LIST_WIDTH  float32 = 300
	SIZE_LIST_HEIGHT float32 = 39.0665647
)

type GUI interface {
	SetMainWindow()
	ShowGUI()
}

// Keeps inner workings of the gui private but allows passing of information from pane to pane
type guiInternal interface {
	startNewSearch([][]string, *util.SearchSettings, time.Time)
}

type gui struct {
	application fyne.App
	mainWindow  fyne.Window

	// Main menu
	mainMenu mainMenu

	// Default view panes
	searchConfigurationPane searchConfigurationPane
	searchResultsPane       searchResultsPane
	corpusPane              corpusPane
}

// Package initialisation just covers the version for now
func init() {
	VERSION = VERSION_MAJOR + "." + VERSION_MINOR + "." + VERSION_PATCH
}

func NewGUI() GUI {
	g := &gui{
		application: app.NewWithID(APP_ID),
	}
	g.mainWindow = g.application.NewWindow(fmt.Sprintf("%s - v%s", MAIN_WINDOW_TITLE, VERSION))

	g.createMainMenu()

	// Always start with default view loaded
	g.createDefaultView()

	return g
}

func (g *gui) createMainMenu() {
	g.mainMenu = newMainMenu(g)
	g.mainWindow.SetMainMenu(g.mainMenu.getMainMenu())
}

func (g *gui) createDefaultView() {
	// Set up individual panes
	g.searchConfigurationPane = newSearchConfigurationPane(g)
	g.searchResultsPane = newSearchResultsPane()
	g.corpusPane = newCorpusPane()

	// Create and set the layout of the panes
	g.mainWindow.SetContent(
		container.NewHSplit(
			container.NewVSplit(
				g.searchConfigurationPane.GetPane(),
				g.searchResultsPane.GetPane(),
			),
			g.corpusPane.GetPane(),
		),
	)
}

func (g *gui) SetMainWindow() {
	g.mainWindow.SetMaster()
}

func (g *gui) ShowGUI() {
	g.mainWindow.ShowAndRun()
}

// Run a new search with the given search terms
func (g *gui) startNewSearch(terms [][]string, searchSettings *util.SearchSettings, startTime time.Time) {
	if err := util.ValidateTypes(terms); err != nil {
		g.searchResultsPane.setResultInfo("Search error:", err.Error())
		return
	}
	g.searchResultsPane.startNewSearch(terms)

	// Run the search
	searchResults, totalInstances, err := util.SearchFor(
		terms,
		searchSettings,
	)
	if err != nil {
		g.searchResultsPane.setResultInfo("Search error:", err.Error())
		return
	}

	searchTime := time.Since(startTime)

	g.searchResultsPane.updateResults(searchResults, totalInstances, searchSettings, startTime, searchTime)

	// TODO: also disable all other widgets associated with search configuration while search is running
}
