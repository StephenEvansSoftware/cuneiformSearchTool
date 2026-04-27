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
type searchResultsPane interface {
	paneGetter

	// Displays the text given directly
	simpleUpdate(markdown []string)
	// Displays a full result
	updateResults(overwrite bool, results []*util.TextSearchResult) // TODO: may change TextSearchResult, update if needed
}

// Contains all objects related to the search result pane
type srpStruct struct {
	// Main container
	//	mainContainer *fyne.Container
	mainContainer fyne.CanvasObject

	// Result container
	results          *fyne.Container
	resultsLabel     *widget.Label
	resultsInfo      *widget.RichText
	resultsAccordion *widget.Accordion

	// Text display container
	textDisplay textDisplayPane
}

func newSearchResultsPane() searchResultsPane {
	srPane := &srpStruct{}

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

func (srp *srpStruct) setupResults() {
	srp.resultsLabel = widget.NewLabel(labelTextResults)
	srp.resultsInfo = widget.NewRichTextFromMarkdown(richTextResultsInfoDefault)
	srp.resultsAccordion = widget.NewAccordion() // TODO: this will likely be taken out into its own file to make scrolling with the interior list possible
	srp.resultsAccordion.MultiOpen = true

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

func (srp *srpStruct) setupTextDisplay() {
	srp.textDisplay = newTextDisplay()
}

func (srp *srpStruct) GetPane() fyne.CanvasObject {
	return srp.mainContainer
}

func (srp *srpStruct) simpleUpdate(markdown []string) {
	for _, line := range markdown {
		srp.resultsInfo.AppendMarkdown(line)
	}
}

func (srp *srpStruct) updateResults(overwrite bool, results []*util.TextSearchResult) {
	// TODO
}
