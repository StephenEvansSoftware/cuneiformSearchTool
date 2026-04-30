package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"cuneiformSearchTool/util"
)

// This file creates the search configuration pane of the gui

const (
	labelTextSearchSettings = "Search:"
	labelTextTerms          = "Search terms:"

	buttonTextAddType = "Add type"
	buttonTextSearch  = "SEARCH"
	buttonTextAddTerm = "Add another\nsearch term"
)

var (
	// TODO: consider enum or mapping here
	confTextSearchScope = []string{
		"All Texts",
		"Select Corpus(s)",
		"Select Text(s)",
	}

	// TODO: consider enum or mapping here
	confTextSearchSettings = []string{
		"Wrap line text during search",
		"Print text meta information",
		"Use POS instead of EPOS",
		"Use FORM instead of FRAG",
	}
)

// Publically facing functionality of the search configuration pane
type searchConfigurationPane interface {
	paneGetter
}

// Contains all objects related to the search configuration pane
type scpStruct struct {
	// Main container
	mainContainer *fyne.Container

	// Search settings label
	labelSearchSettings *widget.Label

	// Horizontal configuration container
	configuration       *fyne.Container
	confScope           *widget.RadioGroup
	confCorpusSelection *widget.CheckGroup
	confSettings        *widget.CheckGroup

	// Search terms label
	labelSearchTerms *widget.Label

	// Horizontal search terms container
	terms *fyne.Container

	// Button to add additional search terms
	termsAdd *widget.Button

	// Button to actually run the search
	searchButton *widget.Button

	// Reference to parent gui object to pass information to
	gui guiInternal
}

func newSearchConfigurationPane(gui guiInternal) searchConfigurationPane {
	scPane := &scpStruct{}

	// Setup main pane components
	scPane.setupConfiguration()
	scPane.setupSearchTerms()
	scPane.setupSearchButton()

	// Setup labels
	scPane.labelSearchSettings = widget.NewLabel(labelTextSearchSettings)
	scPane.labelSearchTerms = widget.NewLabel(labelTextTerms)

	// Create the finalised search configuration pane
	scPane.mainContainer = container.NewVBox(
		scPane.labelSearchSettings, // Settings label
		scPane.configuration,       // Search configuration
		widget.NewSeparator(),
		scPane.labelSearchTerms, // Terms label
		scPane.terms,            // Search terms
		scPane.termsAdd,         // Add additional search terms
		scPane.searchButton,     // Button to actually run the search
	)

	// Save reference to parent to start search
	scPane.gui = gui

	return scPane
}

func (scp *scpStruct) setupConfiguration() {
	// Set up corpus selection so we can enable / disable it depending on selected scope
	scp.confCorpusSelection = widget.NewCheckGroup([]string{}, func(changed []string) {})
	for _, corpusText := range util.ListCorpusAndTexts() { // TODO: expand to select texts in accordion fashion
		scp.confCorpusSelection.Append(corpusText.CorpusName)
	}

	// Radio selection of search scope
	scp.confScope = widget.NewRadioGroup(
		confTextSearchScope,
		func(s string) {
			switch s {
			case confTextSearchScope[0]: // TODO: manually checking position 0 is messy here
				scp.confCorpusSelection.Disable()
			default:
				scp.confCorpusSelection.Enable()
			}
		},
	)
	// Default scope selection of select corpus
	scp.confScope.SetSelected(confTextSearchScope[1]) // TODO: manually checking position 1 is messy here
	// Ensure scope is always set
	scp.confScope.Required = true

	// Search settings
	scp.confSettings = widget.NewCheckGroup(
		confTextSearchSettings,
		func(changed []string) {}, // TODO: is this needed further down the line?
	)

	// Add the search configuration panes in correct order
	scp.configuration = container.NewHBox(
		scp.confScope,
		widget.NewSeparator(),
		container.NewVScroll(scp.confCorpusSelection),
		widget.NewSeparator(),
		scp.confSettings,
	)
}

func (scp *scpStruct) setupSearchTerms() {
	scp.terms = container.NewHBox()
	scp.terms.Add(scp.createNewSearchTerm())
	scp.termsAdd = widget.NewButton(buttonTextAddTerm, func() {
		scp.terms.Add(scp.createNewSearchTerm())
	})
}

func (scp *scpStruct) createNewSearchTerm() *fyne.Container {
	// Overarching search term
	encapsulatedTerm := container.NewVBox()
	// Individual types of the search term
	searchTypes := container.NewVBox()

	// Function to add additional types to the search term
	createType := func() fyne.CanvasObject {
		// Selection for type as well as removal of type
		newType := container.NewHBox()
		newType.Add(widget.NewSelect(append([]string{util.AllProperNounsType}, util.ValidTypes...), func(s string) {}))
		// Button to remove this type from the individual search term
		removeType := widget.NewButton("-", func() {
			// Check if this is the last type in this search term
			if len(searchTypes.Objects) == 1 {
				// Check if this is the last remaining term in the search
				if len(scp.terms.Objects) == 1 {
					// If so, do not remove this type and move on
					// TODO: show an error at this point, or maybe disable the removal if only one type in one term
					return
				}
			}
			// If not, remove this type
			searchTypes.Remove(newType)
			// Check if this was the last type removed from the term
			if len(searchTypes.Objects) == 0 {
				// If so, remove this term
				scp.terms.Remove(encapsulatedTerm)
			}
		})
		newType.Add(removeType)
		return newType
	}

	// Each term must have one initial type so add it now
	searchTypes.Add(createType())

	// Button to add further types to this term
	addType := widget.NewButton(buttonTextAddType, func() {
		searchTypes.Add(createType())
	})

	// Finally, add the list of types along with the button to add more
	encapsulatedTerm.Add(searchTypes)
	encapsulatedTerm.Add(addType)

	return encapsulatedTerm
}

func (scp *scpStruct) setupSearchButton() {
	scp.searchButton = widget.NewButton(buttonTextSearch, func() {
		// TODO TODO TODO: the ui doesn't change until the search if fully complete
		//					so disable the button and set up a go routine to run the function
		//					before re-enabling after the search is complete
		//					also disable all other search settings while search running

		// Run a timer to see how long each search takes
		timeStart := time.Now()

		// Retrieve all search terms
		collectedSearchTerms := [][]string{}
		for sti, st := range scp.terms.Objects { // Overarching search terms container
			// Create a new set of terms
			collectedSearchTerms = append(collectedSearchTerms, []string{})
			for _, sto := range st.(*fyne.Container).Objects { // Types and buttons, and button
				// Check for *fyne.Container as this holds our actual types
				switch t := sto.(type) {
				case *fyne.Container:
					for _, to := range t.Objects { // Types and buttons
						for _, too := range to.(*fyne.Container).Objects {
							switch toot := too.(type) {
							case *widget.Select: // The actual search type we're looking for
								// Ignore selections that haven't been selected
								if selected := toot.Selected; selected != "" {
									// Add the type to our current term
									collectedSearchTerms[sti] = append(collectedSearchTerms[sti], selected)
								}
							case *widget.Button: // "Remove type" button, ignore
							default: // Unknown type, should never happen
								// TODO: handle better?
								panic(fmt.Sprintf("term contains non-expected type: %v", toot))
							}
						}
					}
				case *widget.Button: // "Add type" button, ignore
				default: // Unknown type, should never happen
					// TODO: handle better?
					panic(fmt.Sprintf("terms contain non-expected type: %v", sto))
				}
			}
		}

		// Collect all search configuration and settings
		searchSettings := &util.SearchSettings{}
		switch scp.confScope.Selected { // TODO: have this be handled by changing search settings when clicked instead of getting it here
		case confTextSearchScope[0]: // Search all texts
			searchSettings.SearchAllTexts = true
		case confTextSearchScope[1]: // Search selected corpus'
			searchSettings.SelectedCorpus = scp.confCorpusSelection.Selected
		case confTextSearchScope[2]: // Search selected texts // TODO
			searchSettings.SelectedCorpus = scp.confCorpusSelection.Selected
		default: // Should never happen
			panic("cannot determine configuration search scope")
		}
		// TODO: same as above, have it automatically set the settings whenever clicked
		for _, ss := range scp.confSettings.Selected {
			if ss == confTextSearchSettings[0] {
				searchSettings.SearchWrappedText = true
				continue
			}
			if ss == confTextSearchSettings[1] {
				searchSettings.PrintTextMetaInformation = true
				continue
			}
			if ss == confTextSearchSettings[2] {
				searchSettings.UsePOS = true
				continue
			}
			if ss == confTextSearchSettings[3] {
				searchSettings.UseFORM = true
				continue
			}
		}

		// Inform gui that a new search is happening and let it handle the rest
		scp.gui.startNewSearch(collectedSearchTerms, searchSettings, timeStart)
	})
}

func (scp *scpStruct) getSearchSettings() (*util.SearchSettings, error) {
	return nil, nil
}

// Interface implementations

func (scp *scpStruct) getPane() fyne.CanvasObject {
	return scp.mainContainer
}
