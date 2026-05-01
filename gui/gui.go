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

	// Quick workaround for "tabbing" in markdown
	richTextTab = "&nbsp&nbsp&nbsp&nbsp"
)

type GUI interface {
	SetMainWindow()
	ShowGUI()
}

// Keeps inner workings of the gui private but allows passing of information from pane to pane
type guiInternal interface {
	startNewTermsSearch([][]string, *util.SearchSettings, time.Time)
	startNewMassSearch(*util.SearchSettings, *util.MassSearchSettings, time.Time)
	setViewSearch()
	toggleCorpusPane()
}

type gui struct {
	application fyne.App
	mainWindow  fyne.Window

	// Main menu
	mainMenu mainMenu

	// Search view panes
	searchView              fyne.CanvasObject
	searchConfigurationPane searchConfigurationPane
	searchResultsPane       searchResultsPane

	// Generic panes used across multiple views
	corpusPane corpusPane
	// Track if we want to have corpus visible or not
	corpusView bool
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

	// Create all view panes
	g.createViewPanes()

	// Always start with default view loaded
	g.setViewDefault()

	return g
}

func (g *gui) createMainMenu() {
	g.mainMenu = newMainMenu(g)
	g.mainWindow.SetMainMenu(g.mainMenu.getMainMenu())
}

func (g *gui) createViewPanes() {
	// Generic panes used in multiple views
	g.corpusPane = newCorpusPane()
	g.corpusView = true

	// Search panes
	g.searchConfigurationPane = newSearchConfigurationPane(g)
	g.searchResultsPane = newSearchResultsPane()
	g.searchView = container.NewHSplit(
		container.NewVSplit(
			g.searchConfigurationPane.getPane(),
			g.searchResultsPane.getPane(),
		),
		g.corpusPane.getPane(),
	)
}

func (g *gui) setViewDefault() {
	// Default view is just the search view
	g.setViewSearch()
}

func (g *gui) setViewSearch() {
	g.mainWindow.SetContent(g.searchView)
}

func (g *gui) toggleCorpusPane() {
	g.corpusView = !g.corpusView
	if g.corpusView {
		g.corpusPane.getPane().Show()
	} else {
		g.corpusPane.getPane().Hide()
	}
}

func (g *gui) SetMainWindow() {
	g.mainWindow.SetMaster()
}

func (g *gui) ShowGUI() {
	g.mainWindow.ShowAndRun()
}

// Run a new search with the given search terms
func (g *gui) startNewTermsSearch(terms [][]string, searchSettings *util.SearchSettings, startTime time.Time) {
	if err := util.ValidateTypes(terms); err != nil {
		g.searchResultsPane.setResultInfo("Search error:", err.Error())
		g.searchConfigurationPane.finishSearch()
		return
	}
	// TODO: update and refresh search result pane before search continues
	g.searchResultsPane.startNewSearch(terms)

	// Run the search
	searchResults, totalInstances, err := util.SearchFor(
		terms,
		searchSettings,
	)
	if err != nil {
		g.searchResultsPane.setResultInfo("Search error:", err.Error())
		g.searchConfigurationPane.finishSearch()
		return
	}

	searchTime := time.Since(startTime)

	g.searchResultsPane.updateResults(searchResults, totalInstances, searchSettings, startTime, searchTime)
	g.searchConfigurationPane.finishSearch()

	// TODO: also disable all other widgets associated with search configuration while search is running
}

func (g *gui) startNewMassSearch(searchSettings *util.SearchSettings, massSearchSettings *util.MassSearchSettings, startTime time.Time) {
	util.RunMassSearch(searchSettings, massSearchSettings)
}
