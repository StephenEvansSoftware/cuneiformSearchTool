package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"cuneiformSearchTool/util"
)

var (
	VERSION_MAJOR, VERSION_MINOR, VERSION_PATCH string
	VERSION                                     string
)

const (
	APP_ID            = "cunfeiform-search-tool"
	MAIN_WINDOW_TITLE = "Cuneiform Search Tool"

	// TODO: hacky way of forcing the AccordionItem to show all elements in list while scrolling Accordion
	//	ITEM_SIZE_HEIGHT = 38
	// TODO TODO: remove
	ITEM_SIZE_WIDTH  float32 = 300
	ITEM_SIZE_HEIGHT float32 = 39.0665647

	// TODO: hack
	SIZE_LIST_WIDTH  float32 = 300
	SIZE_LIST_HEIGHT float32 = 39.0665647
)

type GUI interface {
	// TODO: add only things that the main function cares about here
	SetMainWindow()
	ShowGUI()
}

// Keeps inner workings of the gui private but allows passing of information from pane to pane
type guiInternal interface {
	updateResults(string) // TODO: just for testing, will change later
}

type gui struct {
	/*

				LP						RP
			|-------------------------------|
			|	LPT					|	RPT	|
			|-----------------------|		|
			|	LPBL	|	LPBR	|	RPB	|
			|-------------------------------|

			LP = Left Pane			|	Holds the left split in the main window
		>	LPT = LP Top			|	Search input
		>	LPBL = LP Bottom Left	|	List of search results
		>	LPBR = LP Bottom Right	|	Search result text with highlighted occurences
			RP = Right Pane			|	Holds the right split in the main window
			RPT = RP Top			|	List of corpus' and texts therein
			RPB = RP Bottom			|	View of selected text


	*/
	Application fyne.App
	MainWindow  fyne.Window

	LeftPane                 fyne.CanvasObject
	LeftPaneTop              fyne.CanvasObject
	LeftPaneBottom           fyne.CanvasObject
	LeftPaneBottomLeftTop    *widget.RichText // TODO TODO TODO: keep in mind this is markdown
	LeftPaneBottomLeftBottom *widget.Accordion
	LeftPaneBottomRight      *widget.RichText

	RightPane       fyne.CanvasObject
	RightPaneTop    fyne.CanvasObject
	RightPaneBottom fyne.CanvasObject

	// TODO TODO: update this
	searchScope           *widget.RadioGroup
	searchCorpusSelection *widget.CheckGroup
	searchSettings        *widget.CheckGroup

	// TODO TODO TODO: end of legacy, remove when new working

	// Default view panes
	searchConfigurationPane SearchConfigurationPane
	searchResultsPane       SearchResultsPane
	corpusPane              CorpusPane
}

// Package initialisation just covers the version for now
func init() {
	VERSION = VERSION_MAJOR + "." + VERSION_MINOR + "." + VERSION_PATCH
}

func NewGUI() GUI {
	g := &gui{
		Application: app.NewWithID(APP_ID),
	}
	g.MainWindow = g.Application.NewWindow(fmt.Sprintf("%s - v%s", MAIN_WINDOW_TITLE, VERSION))

	g.createDefaultView()

	return g

	/*
		g := &gui{
			Application: app.NewWithID("cuneiform-search-tool"),
		}
		g.MainWindow = g.Application.NewWindow("Cuneiform Search Tool - v" + version) // TODO: have version number pulled from elsewhere
		g.setupMenu()
		g.MainWindow.Resize(fyne.NewSize(float32(800), float32(600)))
		g.setupContent()
		return g, nil // TODO: just setting up panes for now
	*/
}

func (g *gui) createDefaultView() {
	g.searchConfigurationPane = newSearchConfigurationPane(g)
	g.searchResultsPane = newSearchResultsPane()
	g.corpusPane = newCorpusPane()

	g.MainWindow.SetContent(
		container.NewHSplit(
			container.NewVSplit(
				g.searchConfigurationPane.GetPane(),
				g.searchResultsPane.GetPane(),
			),
			g.corpusPane.GetPane(),
		),
	)

}

func (g *gui) setupMenu() {
	// TODO: not needed just yet, panes are the main focus for now
}

func (g *gui) setupContent() {
	g.LeftPane = g.setupLeftPane()
	g.RightPane = g.setupRightPane()
	g.MainWindow.SetContent(container.NewHSplit(g.LeftPane, g.RightPane))
}

func (g *gui) setupLeftPane() fyne.CanvasObject {
	g.LeftPaneTop = g.setupSearchPane()
	searchResultsLabel := widget.NewLabel("Search Results:")
	g.LeftPaneBottomLeftTop = widget.NewRichTextWithText("Search results will be displayed here after initial search")
	g.LeftPaneBottomLeftBottom = widget.NewAccordion()
	g.LeftPaneBottomLeftBottom.MultiOpen = true
	leftPaneBottom := container.NewVScroll(container.NewVBox(searchResultsLabel, g.LeftPaneBottomLeftTop, g.LeftPaneBottomLeftBottom))
	g.LeftPaneBottomRight = widget.NewRichTextWithText("Text of result here when selected")
	g.LeftPaneBottom = container.NewHSplit(leftPaneBottom, container.NewScroll(g.LeftPaneBottomRight)) // TODO
	return container.NewVSplit(g.LeftPaneTop, g.LeftPaneBottom)
}

func (g *gui) setupRightPane() fyne.CanvasObject {
	// Right Top
	// Get list of all corpus and texts to display
	corpusTexts := util.ListCorpusAndTexts()
	corpusAcc := widget.NewAccordion()
	corpusAcc.MultiOpen = true
	for _, corpusText := range corpusTexts {
		textListLen := len(corpusText.Texts)
		textList := widget.NewList(
			func() int { return textListLen },
			func() fyne.CanvasObject { return widget.NewLabel("Template") },
			func(liID widget.ListItemID, co fyne.CanvasObject) {
				co.(*widget.Label).SetText(corpusText.Texts[liID].TextName)
			},
		)
		sizeW := ITEM_SIZE_WIDTH
		sizeH := float32(textListLen) * ITEM_SIZE_HEIGHT
		textListSize := fyne.NewSize(sizeW, sizeH)
		textListGridWrap := container.NewGridWrap(textListSize, textList)
		corpusAccItem := widget.NewAccordionItem(corpusText.CorpusName, container.NewHScroll(textListGridWrap))
		//		corpusAccItem := widget.NewAccordionItem(corpusText.CorpusName, textList)
		corpusAcc.Append(corpusAccItem)
		textList.OnSelected = func(id widget.ListItemID) {
			//			addW := float32(1000)
			//			addH := float32(1000)
			//			textList.BaseWidget.Size().AddWidthHeight(addW, addH)
			//			textList.Resize(textList.Size().AddWidthHeight(1000,1000))
			//			corpusAcc.Resize(corpusAcc.Size().AddWidthHeight(1000,1000))
			//			g.RightPaneTop.Resize(g.RightPaneTop.Size().AddWidthHeight(1000,1000))
			//corpusAcc.BaseWidget.Size().AddWidthHeight(addW, addH)

		}
	}
	//	g.RightPaneTop = container.NewBorder(widget.NewLabel("Sources:"), nil, nil, nil, container.NewScroll(corpusAcc))
	//	g.RightPaneTop = container.NewVBox(widget.NewLabel("Sources:"), corpusAcc)
	g.RightPaneTop = container.NewScroll(corpusAcc)

	// Right Bottom
	// TODO: show selected text from right top
	rpbcl := widget.NewLabelWithStyle("Loaded Text:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	rpbct := widget.NewLabelWithStyle("Q00TEST", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	rpbc := container.NewBorder(nil, nil, rpbcl, nil, rpbct)
	rpbcrt := widget.NewRichText()
	rpbcrt.AppendMarkdown("test simple string")
	rpbcrt.AppendMarkdown("# header")
	rpbcrt.AppendMarkdown("- item")
	rpbcrt.AppendMarkdown("- item2")
	rpbcrt.AppendMarkdown("1. numbered")
	rpbcrt.AppendMarkdown("1. numbered")

	rpb := container.NewBorder(
		rpbc,
		nil,
		nil,
		nil,
		rpbcrt,
	)
	g.RightPaneBottom = container.NewScroll(rpb)
	return container.NewVSplit(g.RightPaneTop, g.RightPaneBottom)
}

func (g *gui) setupSearchPane() fyne.CanvasObject {
	/*
		_______________________________________________
		|Search:		__________	Search Settings:
		|[]	All			|-C1 []	 |		[] Line Wrap
		|[] Corpus(s)	|-C2 []	 |
		|[] Text(s)		|-C3 []	 |
		|				|>Q001 []|
		|				|>Q002 []|
		|				|________|
		|Search Terms:
		| [N v]	[N v]-
		|	+	[V v]-	 +
		|		  +
		|_______________________________________________

		Search:
			top left pane
			has:
				label
				exclusive checkboxes
				pane for corpus and texts:
					ALL: greyed out, disabled
					CORPUS: only corpus selectable
					TEXT: add texts to accordion
						note: if corpus selected, selects all texts underneath
						can then untick whichever ones you want

			top right pane
			has:
				label
				non-exlucive checkboxes

			bottom pane
			has:
				label
				scrollable <--> for input search terms
					terms have dropdown
					- & + buttons
						- removes that term
						+ adds another one for that position
				+ button
					adds more terms

	*/

	// Overarching pane for entire search settings, terms, and user functionality
	searchPane := container.NewVBox()

	// Label at the top
	searchLabel := widget.NewLabel("Search:")
	searchPane.Add(searchLabel)

	// Radio selection, manual selection, and settings displayed horizontally
	searchConfiguration := container.NewHBox()

	// Setting up manual selection here so we have a reference to enable / disable for the radio group below
	// Manual selection of search scope if not "all" is selected"
	corpusSelectionManual := widget.NewCheckGroup([]string{}, func(changed []string) {})
	for _, corpusText := range util.ListCorpusAndTexts() {
		corpusSelectionManual.Append(corpusText.CorpusName)
	}
	g.searchCorpusSelection = corpusSelectionManual

	// Radio selection of search scope
	corpusSelectionRadio := widget.NewRadioGroup(
		[]string{"All Texts", "Select Corpus(s)", "Select Text(s)"},
		func(s string) {
			if s == "Select Corpus(s)" {
				g.searchCorpusSelection.Enable()
			} else {
				g.searchCorpusSelection.Disable()
			}
		},
	)
	corpusSelectionRadio.SetSelected("All Texts")
	corpusSelectionRadio.Required = true
	g.searchScope = corpusSelectionRadio
	searchConfiguration.Add(g.searchScope)

	// Add the manual selection now, after the radio selection has been added
	searchConfiguration.Add(container.NewVScroll(g.searchCorpusSelection))

	// Search settings
	searchSettings := widget.NewCheckGroup([]string{}, func(changed []string) {})
	// TODO TODO TODO: redo this atrocity of literal strings
	searchSettings.Append("Search over multiple lines")
	searchSettings.Append("Print text meta information")
	searchSettings.Append("Use POS instead of EPOS")
	searchSettings.Append("Use FORM instead of FRAG")
	g.searchSettings = searchSettings
	searchConfiguration.Add(container.NewVScroll(g.searchSettings))

	// Add the finalised horizontal search configuration
	searchPane.Add(searchConfiguration)

	// Search terms
	searchPane.Add(widget.NewLabel("Terms:"))
	searchTerms := container.NewHBox()
	searchTerms.Add(g.createNewSearchTerm(searchTerms))
	searchPane.Add(container.NewHScroll(searchTerms))
	searchPane.Add(widget.NewButton("Add another\nsearch term", func() {
		searchTerms.Add(g.createNewSearchTerm(searchTerms))
	}))

	// Finally, the search button
	searchButton := widget.NewButton("SEARCH", func() {})
	searchButton.OnTapped = func() {
		// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO: this needs to be put into a go routine and disable ALL search configuration including buttons etc until the search has been completed. reason for go routine is that the ui does not updated anything until the function has 100% completed, even with Refresh() calls

		// Run a timer to see how long each search takes
		timeStart := time.Now()

		// Ensure user cannot repeatedly click the search button while a search is active
		searchButton.Disable()
		searchButton.Importance = widget.LowImportance
		//searchButton.Refresh() // TODO: not working...

		// Retrieve all search terms
		collectedSearchTerms := [][]string{}
		for i, st := range searchTerms.Objects { // Overarching search term container
			// Create a new set of terms
			collectedSearchTerms = append(collectedSearchTerms, []string{})
			for _, sto := range st.(*fyne.Container).Objects { // Individual elements within overarching container (terms, buttons, etc)
				// Check for *fyne.Container as this is the object containing our terms
				switch t := sto.(type) {
				case *fyne.Container:
					for _, to := range t.Objects { // This are the actual search terms (term + add more terms)
						for _, too := range to.(*fyne.Container).Objects { // Individual term and button to remove
							switch toot := too.(type) {
							case *widget.Select: // Our actual search term is found here
								if selected := toot.Selected; selected != "" {
									// Add the term to our running total
									collectedSearchTerms[i] = append(collectedSearchTerms[i], selected)
								}
							case *widget.Button: // Do nothing, this is the "remove type" button
							default: // Unknown type, should never happen
								panic(fmt.Sprintf("searchTerm term contains a non-expected type: %s", toot)) // TODO: remove if all testing is fine
							}
						}
					}
				case *widget.Button: // Do nothing, this is the "add another type" button
				default: // Unknown type, should never happen
					panic(fmt.Sprintf("searchTerm container contains a non-expected type: %s", sto)) // TODO: remove if all testing is fine
				}
			}
		}
		// Reset previous search results and show what is currently being searched for
		g.LeftPaneBottomLeftTop.ParseMarkdown(fmt.Sprintf("Searching for: %+v\n", collectedSearchTerms))
		g.LeftPaneBottomLeftTop.Refresh() // TODO: not working...
		g.RunSearch(collectedSearchTerms, timeStart)
		searchButton.Importance = widget.MediumImportance
		searchButton.Enable()
		searchButton.Refresh()

	}
	searchPane.Add(searchButton)

	return searchPane
}

func (g *gui) createNewSearchTerm(parent *fyne.Container) *fyne.Container {
	encapsulatedTerm := container.NewVBox()
	searchTerms := container.NewVBox()

	individualTermCreate := func() fyne.CanvasObject {
		searchTerm := container.NewHBox()
		searchTerm.Add(
			widget.NewSelect(util.ValidTypes, func(s string) {}),
		)
		removeTerm := widget.NewButton("-", func() {
			// Check if this is the last type in this search term
			if len(searchTerms.Objects) == 1 {
				// Check if this is the last remaining term in the search
				if len(parent.Objects) == 1 {
					// If so, do not remove this and simply move on
					return
				}
			}
			// If not, remove and check
			searchTerms.Remove(searchTerm)
			if len(searchTerms.Objects) == 0 {
				parent.Remove(encapsulatedTerm)
			}
		})
		searchTerm.Add(removeTerm)
		return searchTerm
	}
	searchTerms.Add(individualTermCreate())

	addMoreTermsButton := widget.NewButton("Add more types", func() {
		searchTerms.Add(
			individualTermCreate(),
		)
	})
	encapsulatedTerm.Add(searchTerms)
	encapsulatedTerm.Add(addMoreTermsButton)
	return encapsulatedTerm
}

// GUI interface implementations for gui struct
func (g *gui) SetMainWindow() {
	g.MainWindow.SetMaster()
}

func (g *gui) ShowGUI() {
	g.MainWindow.ShowAndRun()
}

// Search functionality here
// TODO: remove from gui and put in own package, this is just quick and dirty to get it running asap
func (g *gui) RunSearch(terms [][]string, timeStart time.Time) {
	/*
		This will do the following:
		1. Validate the search terms
		2. Run the search
		3. Display which texts and which lines the search results was found on
	*/

	if err := util.ValidateTypes(terms); err != nil {
		g.LeftPaneBottomLeftTop.ParseMarkdown(fmt.Sprintf("SEARCH ERROR: %s", err.Error()))
		return
	}

	// TODO TODO TODO: apss in all other options selected from search configuration
	scopeAllTexts := g.searchScope.Selected == "All Texts" // TODO: have the selections be vars we can check instead of checking raw strings every time
	searchSettings := &util.SearchSettings{}
	for _, sss := range g.searchSettings.Selected {
		// TODO TODO TODO: change to a switch or a map or something
		if sss == "Search over multiple lines" {
			searchSettings.SearchWrappedText = true
		}
		if sss == "Print text meta information" {
			searchSettings.PrintTextMetaInformation = true
		}
		if sss == "Use POS instead of EPOS" {
			searchSettings.UsePOS = true
		}
		if sss == "Use FORM instead of FRAG" {
			searchSettings.UseFORM = true
		}
	}
	searchResults, totalInstances, err := util.SearchFor(
		terms,
		scopeAllTexts,
		g.searchCorpusSelection.Selected,
		[]string{}, // TODO: not searching individual texts yet
		searchSettings,
	)
	if err != nil {
		g.LeftPaneBottomLeftTop.ParseMarkdown(fmt.Sprintf("SEARCH ERROR: %s", err.Error()))
		return
	}
	g.LeftPaneBottomLeftTop.AppendMarkdown(fmt.Sprintf("Duration of search: %s\n", time.Since(timeStart)))
	g.LeftPaneBottomLeftTop.AppendMarkdown(fmt.Sprintf(
		"Texts containing search terms: %d\n",
		len(searchResults),
	))
	g.LeftPaneBottomLeftTop.AppendMarkdown(fmt.Sprintf(
		"Total instances found across all texts: %d\n",
		totalInstances,
	))
	g.LeftPaneBottomLeftBottom.Items = []*widget.AccordionItem{}
	for _, sr := range searchResults {

		listLen := len(sr.Results)
		textResultsList := widget.NewList(
			func() int { return listLen },
			func() fyne.CanvasObject { return widget.NewLabel("tmp") },
			func(lID widget.ListItemID, co fyne.CanvasObject) {
				co.(*widget.Label).SetText(fmt.Sprintf("Line: %d, Position: %d, Last Entry:%s", sr.Results[lID].OverallLineIndex, sr.Results[lID].Position, sr.Results[lID].LastEntry))
			},
		)

		accSizeW := ITEM_SIZE_WIDTH
		accSizeH := float32(listLen) * ITEM_SIZE_HEIGHT
		resultAccSize := fyne.NewSize(accSizeW, accSizeH)
		resultGridWrap := container.NewGridWrap(resultAccSize, textResultsList)

		rhsgw := container.NewHScroll(resultGridWrap)

		resultAccordionItem := widget.NewAccordionItem(fmt.Sprintf("%s: %d results", sr.TextName, listLen), rhsgw)
		//		resultAccordionItem := widget.NewAccordionItem(fmt.Sprintf("%s: %d results", sr.TextName, listLen), resultGridWrap)
		//		resultAccordionItem := widget.NewAccordionItem(fmt.Sprintf("%s: %d results", sr.TextName, listLen), textResultsList)
		//		g.LeftPaneBottomLeftBottom.Append(resultGridWrap)
		g.LeftPaneBottomLeftBottom.Append(resultAccordionItem)

		// Set up the text loading on selection of a result
		textResultsList.OnSelected = func(id widget.ListItemID) {
			printableText, err := util.GetPrintableText(sr.TextPath, searchSettings.PrintTextMetaInformation, searchSettings.UsePOS, searchSettings.UseFORM)
			if err != nil {
				// TODO TODO TODO: handle this
			}
			g.LeftPaneBottomRight.ParseMarkdown(printableText)
		}
	}
}

func TestNewPane() {

	testApp := app.NewWithID("test-app")
	testWindow := testApp.NewWindow("test-app-window")
	testWindow.Resize(fyne.NewSize(float32(800), float32(600)))

	// Search pane
	//	nsp := newSearchConfigurationPane()
	//	testWindow.SetContent(nsp.GetPane())

	// Result pane
	//	nrp := newSearchResultsPane()
	//	testWindow.SetContent(nrp.GetPane())

	// Corpus pane
	ncp := newCorpusPane()
	testWindow.SetContent(ncp.GetPane())

	testWindow.SetMaster()
	testWindow.ShowAndRun()
}

func (g *gui) updateResults(s string) {
	g.searchResultsPane.simpleUpdate([]string{s})
}
