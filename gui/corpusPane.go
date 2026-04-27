package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"cuneiformSearchTool/util"
)

const (
	labelTextCorpus = "Corpus Texts:"
)

type corpusPane interface {
	paneGetter
}

type cpStruct struct {
	mainContainer fyne.CanvasObject

	corpus          *fyne.Container
	corpusLabel     *widget.Label
	corpusAccordion *widget.Accordion

	textDisplay textDisplayPane
}

func newCorpusPane() corpusPane {
	cpp := &cpStruct{}

	cpp.setupCorpus()
	cpp.setupTextDisplay()

	cpp.mainContainer = container.NewVSplit(
		cpp.corpus,
		container.NewScroll(cpp.textDisplay.GetPane()),
	)

	return cpp
}

func (cpp *cpStruct) setupCorpus() {
	cpp.corpusLabel = widget.NewLabel(labelTextCorpus)

	cpp.corpusAccordion = widget.NewAccordion()
	cpp.corpusAccordion.MultiOpen = true

	// Set up the AccordionItems with the corpus text information
	for _, corpusText := range util.ListCorpusAndTexts() {
		// Create a list of texts within the corpus
		corpusTextLen := len(corpusText.Texts)
		corpusTextList := widget.NewList(
			func() int { return corpusTextLen },                        // Length function for the list
			func() fyne.CanvasObject { return widget.NewLabel("tmp") }, // Initial list items
			func(liID widget.ListItemID, co fyne.CanvasObject) {
				co.(*widget.Label).SetText(corpusText.Texts[liID].TextName)
			}, // List update function
		)

		// Set up selection function to load selected text
		corpusTextList.OnSelected = func(liID widget.ListItemID) {
			// TODO TODO TODO: this is still messy as hell, refactor
			printableText, err := util.GetPrintableText(corpusText.Texts[liID].TextPath, true, false, false)
			if err != nil {
				panic(err.Error()) // TODO TODO TODO: 100% change this
			}
			cpp.textDisplay.SetTitle(corpusText.Texts[liID].TextName)
			cpp.textDisplay.ParseMarkdown([]string{printableText})
		}

		// Create custom sized container for the list to avoid fyne imposing MinSize and making it unusable
		// TODO: better way to do this?
		corpusTextListSize := fyne.NewSize(
			SIZE_LIST_WIDTH, // TODO: different width?
			SIZE_LIST_HEIGHT*float32(corpusTextLen), // TODO: have a max?
		)
		corpusTextListContainer := container.NewGridWrap(
			corpusTextListSize,
			corpusTextList,
		)

		// Add the list to the accordion
		corpusTextAccItem := widget.NewAccordionItem(
			corpusText.CorpusName,
			corpusTextListContainer,
		)
		cpp.corpusAccordion.Append(corpusTextAccItem)
	}

	cpp.corpus = container.NewBorder(
		container.NewVBox(
			cpp.corpusLabel,
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		container.NewScroll(cpp.corpusAccordion),
	)
}

func (cpp *cpStruct) setupTextDisplay() {
	cpp.textDisplay = newTextDisplay()
}

func (cpp *cpStruct) GetPane() fyne.CanvasObject {
	return cpp.mainContainer
}
