package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"cuneiformSearchTool/util"
)

const (
	labelTextResults           = "# **Results:**"
	labelTextResultDisplay     = "# **Result Display:**"
	richTextResultsTextDefault = "Run a search to see results here."
)

type searchResultsPane interface {
	paneGetter

	setParentTab(*container.TabItem)
	setError(err string)
	searchStarted(searchSettings *util.SearchSettings, searchTerms [][]string)
	searchFinished(searchResults []*util.CorpusSearchResult, totalInstances int, searchTimeStart time.Time, searchTimeDuration time.Duration)
}

type srpStruct struct {
	mainContainer fyne.CanvasObject

	results          *fyne.Container
	resultsText      *widget.RichText
	resultsAccordion *widget.Accordion

	searchSettings *util.SearchSettings
	searchTerms    [][]string
	corpusScope    string

	resultDisplay textDisplayPane

	parentTab *container.TabItem
	gui       guiInternal
}

func newSearchResultsPane(gui guiInternal) searchResultsPane {
	srp := &srpStruct{gui: gui}

	srp.setupResults()
	srp.setupResultDisplay()

	srp.mainContainer = container.NewBorder(
		nil, nil,
		container.NewHBox(
			container.NewVScroll(
				container.NewVBox(
					widget.NewRichTextFromMarkdown(labelTextResults),
					widget.NewSeparator(),
					srp.results,
				),
			),
			widget.NewSeparator(),
		), nil,
		container.NewScroll(
			container.NewBorder(
				container.NewVBox(
					widget.NewRichTextFromMarkdown(labelTextResultDisplay),
					widget.NewSeparator(),
				), nil, nil, nil,
				srp.resultDisplay.getPane(),
			),
		),
	)

	return srp
}

func (srp *srpStruct) setupResults() {
	srp.resultsText = widget.NewRichTextFromMarkdown(richTextResultsTextDefault)
	srp.resultsAccordion = widget.NewAccordion()
	srp.resultsAccordion.MultiOpen = true
	srp.results = container.NewVBox(
		srp.resultsText,
		srp.resultsAccordion,
	)
}

func (srp *srpStruct) setupResultDisplay() {
	srp.resultDisplay = newTextDisplay()
}

func (srp *srpStruct) setError(err string) {
	srp.resultsText.ParseMarkdown(fmt.Sprintf("Error during search:\n\n%s", err))
}

func (srp *srpStruct) searchStarted(searchSettings *util.SearchSettings, searchTerms [][]string) {
	srp.resultsText.ParseMarkdown("Searching for following terms:")
	srp.resultsText.AppendMarkdown(fmt.Sprintf("%v", searchTerms))

	srp.resultsText.AppendMarkdown("**Search Settings Used:**")
	corpusScope := "All"
	if !searchSettings.SearchAllCorpuses {
		corpusScope = fmt.Sprintf("%v", searchSettings.SelectedCorpus)
	}
	srp.resultsText.AppendMarkdown(fmt.Sprintf("``- Selected Corpuses       : %s``", corpusScope))
	srp.resultsText.AppendMarkdown(fmt.Sprintf("``- Selected Texts          : %s``", "All")) // TODO
	srp.resultsText.AppendMarkdown(fmt.Sprintf("``- Wrap line of text       : %v``", searchSettings.SearchWrappedText))
	srp.resultsText.AppendMarkdown(fmt.Sprintf("``- Ignore overlap          : %v``", searchSettings.IgnoreOverlap))
	srp.resultsText.AppendMarkdown(fmt.Sprintf("``- Use POS instead of EPOS : %v``", searchSettings.UsePOS))
	srp.resultsText.AppendMarkdown(fmt.Sprintf("``- Show FORM over FRAG     : %v``", searchSettings.UseFORM))
	srp.searchSettings = searchSettings
	srp.searchTerms = searchTerms

	srp.corpusScope = corpusScope
}

func (srp *srpStruct) searchFinished(searchResults []*util.CorpusSearchResult, totalInstances int, searchTimeStart time.Time, searchTimeDuration time.Duration) {

	totalTextResults := 0

	for _, cr := range searchResults {
		corpusAcc := widget.NewAccordion()
		corpusAcc.MultiOpen = true
		totalTextResults += len(cr.TextResults)
		for _, tr := range cr.TextResults {
			resultsListLen := len(tr.Results)
			resultsList := widget.NewList(
				func() int { return resultsListLen },
				func() fyne.CanvasObject { return widget.NewLabel("tmp") },
				func(liID widget.ListItemID, co fyne.CanvasObject) {
					co.(*widget.Label).SetText(fmt.Sprintf(
						"    %d : %s : %d : %s",
						tr.Results[liID].OverallLineIndex,
						tr.Results[liID].Label,
						tr.Results[liID].Position,
						tr.Results[liID].LastEntry,
					))
				},
			)
			resultsListSize := fyne.NewSize(
				SIZE_LIST_WIDTH,
				SIZE_LIST_HEIGHT*float32(resultsListLen),
			)
			resultsListContainer := container.NewGridWrap(
				resultsListSize,
				resultsList,
			)

			resultAccItem := widget.NewAccordionItem(
				fmt.Sprintf("%s: %d results", tr.TextName, resultsListLen),
				resultsListContainer,
			)
			corpusAcc.Append(resultAccItem)
			resultsList.OnSelected = func(liID widget.ListItemID) {
				printSettings := srp.gui.getPrintSettings()
				printSettings.UsePos = srp.searchSettings.UsePOS
				printableText, err := util.GetPrintableText(
					tr.TextPath,
					printSettings,
				)
				if err != nil {
					panic(err) // TODO: handle better
				}
				srp.resultDisplay.SetTitle(tr.TextName)
				display := []string{printableText.Highlight(tr.Results[liID].OverallLineIndex, tr.Results[liID].Position, srp.searchSettings.TermsLength)}
				if printSettings.PrintTextMetaInformation {
					display = append([]string{printableText.MetaInformation.Sprint()}, display...)
				}
				srp.resultDisplay.ParseMarkdown(display)
				resultsList.UnselectAll()
				resultsList.Highlight(liID)
			}
		}

		corpusAccLeftPadding := container.NewHBox(corpusAcc)
		corpusAccLeftPadding.Layout = layout.NewCustomPaddedLayout(0, 0, 10, 0)

		corpusAccItem := widget.NewAccordionItem(fmt.Sprintf("%s: %d", cr.CorpusName, len(cr.TextResults)), corpusAccLeftPadding)
		srp.resultsAccordion.Append(corpusAccItem)
	}

	srp.parentTab.Text = fmt.Sprintf("%s : %v", srp.corpusScope, srp.searchTerms)
}

func (srp *srpStruct) setParentTab(tab *container.TabItem) {
	srp.parentTab = tab
}

func (srp *srpStruct) getPane() fyne.CanvasObject {
	return srp.mainContainer
}
