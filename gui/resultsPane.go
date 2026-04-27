package gui

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

    "cuneiformSearchTool/util"
)

// This file creates the search result pane of the gui

const (
	labelTextResults = "Results:"

	richTextResultsInfoDefault = "*Results will be displayed here after initial search.*"
)

// Publically facing functionality of the search results pane
// TODO TODO TODO: make this and other interfaces private as they are only to be used within the gui package
type SearchResultsPane interface {
	PaneGetter

	// Displays the text given directly
	simpleUpdate(markdown []string)
	// Displays a full result
	updateResults(overwrite bool, results []*util.TextSearchResult) // TODO: may change TextSearchResult, update if needed
}

// Contains all objects related to the search result pane
type searchResultsPane struct {
	// Main container
//	mainContainer *fyne.Container
	mainContainer fyne.CanvasObject

	// Result container
	results *fyne.Container
	resultsLabel *widget.Label
	resultsInfo *widget.RichText
	resultsAccordion *widget.Accordion

	// Text display container
	textDisplay TextDisplayPane
}

func newSearchResultsPane() SearchResultsPane {
	srPane := &searchResultsPane{}

	// Setup main pane components
	srPane.setupResults()
	srPane.setupTextDisplay()

	// Create the finalised search result pane
	srPane.mainContainer = container.NewHSplit(
		srPane.results,
		srPane.textDisplay.GetPane(),
	)
	return srPane
}

func (srp *searchResultsPane) setupResults() {
	srp.resultsLabel = widget.NewLabel(labelTextResults)
	srp.resultsInfo = widget.NewRichTextFromMarkdown(richTextResultsInfoDefault)
	srp.resultsAccordion = widget.NewAccordion() // TODO: this will likely be taken out into its own file to make scrolling with the interior list possible
	srp.resultsAccordion.MultiOpen = true


	/*
	srp.results = container.New(
		layout.NewVBoxLayout(),
		srp.resultsLabel,
		widget.NewSeparator(),
		container.NewScroll(srp.resultsInfo),
		srp.resultsAccordion,
	)
	srp.results = container.NewVBox(
		srp.resultsLabel,
		widget.NewSeparator(),
		container.NewScroll(srp.resultsInfo),
		srp.resultsAccordion,
	)
	*/
	//srp.results = container.NewVBox(
	//srp.results = container.NewMax(
	srp.results = container.NewBorder(
		container.NewVBox(
			srp.resultsLabel,
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		container.NewScroll(srp.resultsInfo),
		srp.resultsAccordion,
	)
}


func (srp *searchResultsPane) setupTextDisplay() {
	srp.textDisplay = newTextDisplay()
}

func (srp *searchResultsPane) GetPane() fyne.CanvasObject {
	return srp.mainContainer
}

func (srp *searchResultsPane) simpleUpdate(markdown []string) {
	for _, line := range markdown {
		srp.resultsInfo.AppendMarkdown(line)
	}
}

func (srp *searchResultsPane) updateResults(overwrite bool, results []*util.TextSearchResult) {
	// TODO
}
