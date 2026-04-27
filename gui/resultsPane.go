package gui

import (
	"fmt"

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

	// Indicate a new search is happening
	startNewSearch(markdown string, optional ...string)

	// Displays the text given directly
	simpleUpdate(markdown string, optional ...string)
	// Displays a full result
	updateResults(overwrite bool, results []*util.TextSearchResult, searchSettings *util.SearchSettings) // TODO: may change TextSearchResult, update if needed
	// TODO: not happy with the search settings being passed in here, should have something else for it
}

// Contains all objects related to the search result pane
type srpStruct struct {
	// Main container
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
		container.NewScroll(srPane.textDisplay.GetPane()),
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
		container.NewScroll(container.NewVBox(
			container.NewVScroll(srp.resultsInfo),
			srp.resultsAccordion,
		)),
	)
}

func (srp *srpStruct) setupTextDisplay() {
	srp.textDisplay = newTextDisplay()
}

func (srp *srpStruct) GetPane() fyne.CanvasObject {
	return srp.mainContainer
}

func (srp *srpStruct) startNewSearch(markdown string, optional ...string) {
	// Clear current text display by parsing the markdown given
	srp.resultsInfo.ParseMarkdown(markdown)
	for _, line := range optional { // TODO: this is the same as below, make it a function or such
		srp.resultsInfo.AppendMarkdown(line)
	}
}

func (srp *srpStruct) simpleUpdate(markdown string, optional ...string) {
	srp.resultsInfo.AppendMarkdown(markdown)
	for _, line := range optional {
		srp.resultsInfo.AppendMarkdown(line)
	}
}

func (srp *srpStruct) updateResults(overwrite bool, searchResults []*util.TextSearchResult, searchSettings *util.SearchSettings) {
	// Clear previous results
	srp.resultsAccordion.Items = []*widget.AccordionItem{}
	for _, sr := range searchResults {
		resultsListLen := len(sr.Results)
		resultsList := widget.NewList(
			func() int { return resultsListLen },
			func() fyne.CanvasObject { return widget.NewLabel("tmp") },
			func(liID widget.ListItemID, co fyne.CanvasObject) {
				// TODO: this result display is messy, fix
				co.(*widget.Label).SetText(fmt.Sprintf("Lines: %d, Position: %d, Last Entry: %s", sr.Results[liID].OverallLineIndex, sr.Results[liID].Position, sr.Results[liID].LastEntry))
			},
		)

		// TODO: better way to do this?
		resultsListSize := fyne.NewSize(
			SIZE_LIST_WIDTH,
			SIZE_LIST_HEIGHT*float32(resultsListLen),
		)
		resultsListContainer := container.NewGridWrap(
			resultsListSize,
			resultsList,
		)

		resultAccItem := widget.NewAccordionItem(
			fmt.Sprintf("%s: %d results", sr.TextName, resultsListLen), // TODO: different string here?
			resultsListContainer,
		)
		srp.resultsAccordion.Append(resultAccItem)

		resultsList.OnSelected = func(liID widget.ListItemID) {
			printableText, err := util.GetPrintableText(sr.TextPath, searchSettings.PrintTextMetaInformation, searchSettings.UsePOS, searchSettings.UseFORM)
			if err != nil {
				panic(err.Error()) // TODO TODO TODO: handle this better, just panic to know it's occured for now (if it ever does)
			}
			srp.textDisplay.SetTitle(sr.TextName)
			srp.textDisplay.ParseMarkdown([]string{printableText})
		}
	}
}
