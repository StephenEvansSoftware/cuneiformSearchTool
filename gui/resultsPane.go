package gui

import (
	"fmt"
	"time"

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
	startNewSearch(terms [][]string)

	// Displays the text given directly
	setResultInfo(markdown string, optional ...string)

	// Displays a full result
	updateResults(results []*util.CorpusSearchResult, totalInstances int, searchSettings *util.SearchSettings, startTime time.Time, searchTime time.Duration)
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
			srp.resultsInfo,
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

func (srp *srpStruct) startNewSearch(terms [][]string) {
	// Clear current text display by parsing the markdown given
	srp.resultsInfo.ParseMarkdown("Searching for:")
	for _, term := range terms {
		srp.resultsInfo.AppendMarkdown(fmt.Sprintf("%s%+v", richTextTab, term))
	}
	// Clear accordion of previous results
	srp.resultsAccordion.Items = []*widget.AccordionItem{}
	srp.results.Refresh()
}

// Overwrite anything current in the result info
func (srp *srpStruct) setResultInfo(markdown string, optional ...string) {
	srp.resultsInfo.ParseMarkdown(markdown)
	for _, line := range optional {
		// TODO: always indenting optional strings for now, change if needed
		srp.resultsInfo.AppendMarkdown(fmt.Sprintf("%s%s", richTextTab, line))
	}
	srp.results.Refresh()
}

func (srp *srpStruct) updateResults(searchResults []*util.CorpusSearchResult, totalInstances int, searchSettings *util.SearchSettings, startTime time.Time, searchTime time.Duration) {
	// Clear previous results
	srp.resultsAccordion.Items = []*widget.AccordionItem{}

	// Keep track of total texts containing results to display later
	totalTextResults := 0

	// Range over the corpus results
	for _, cr := range searchResults {
		// Each corpus will be its own accordion holding all text results therein
		corpusAcc := widget.NewAccordion()
		corpusAcc.MultiOpen = true
		totalTextResults += len(cr.TextResults)
		// Range over the search results within the current corpus
		for _, tr := range cr.TextResults {
			resultsListLen := len(tr.Results)
			resultsList := widget.NewList(
				func() int { return resultsListLen },
				func() fyne.CanvasObject { return widget.NewLabel("tmp") },
				func(liID widget.ListItemID, co fyne.CanvasObject) {
					// TODO: this result display is messy, fix
					co.(*widget.Label).SetText(fmt.Sprintf(
						//"Line:%d, Label:%s, Position:%d, Last Entry: %s",
						"%d : %s : %d : %s",
						tr.Results[liID].OverallLineIndex,
						tr.Results[liID].Label,
						tr.Results[liID].Position,
						tr.Results[liID].LastEntry,
					))
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
				fmt.Sprintf("%s: %d results", tr.TextName, resultsListLen), // TODO: different string here?
				resultsListContainer,
			)
			corpusAcc.Append(resultAccItem)

			resultsList.OnSelected = func(liID widget.ListItemID) {
				printableText, err := util.GetPrintableText(
					tr.TextPath,
					searchSettings.PrintTextMetaInformation,
					searchSettings.UsePOS,
					searchSettings.UseFORM,
				)
				if err != nil {
					panic(err.Error()) // TODO TODO TODO: handle this better, just panic to know it's occured for now (if it ever does)
				}
				srp.textDisplay.SetTitle(tr.TextName)
				srp.textDisplay.ParseMarkdown([]string{printableText})
			}
		}
		corpusAccItem := widget.NewAccordionItem(fmt.Sprintf("%s: %d", cr.CorpusName, len(cr.TextResults)), corpusAcc)
		srp.resultsAccordion.Append(corpusAccItem)
	}

	// Update the results information
	srp.resultsInfo.AppendMarkdown(fmt.Sprintf("Search duration: %s\n", searchTime))
	srp.resultsInfo.AppendMarkdown(fmt.Sprintf("Result duration: %s\n", time.Since(startTime)))
	srp.resultsInfo.AppendMarkdown(fmt.Sprintf(
		"Corpus' containing search terms: %d\n",
		len(searchResults),
	))
	srp.resultsInfo.AppendMarkdown(fmt.Sprintf(
		"Texts containing search terms: %d\n",
		totalTextResults,
	))
	srp.resultsInfo.AppendMarkdown(fmt.Sprintf(
		"Total instances found across all texts: %d\n",
		totalInstances,
	))
	srp.resultsInfo.AppendMarkdown("---")
	srp.resultsInfo.AppendMarkdown("Results are shown as:")
	srp.resultsInfo.AppendMarkdown(fmt.Sprintf("%s Line : Label : Position : LastEntry", richTextTab))

	srp.results.Refresh()
}
