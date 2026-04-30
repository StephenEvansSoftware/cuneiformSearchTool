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
	labelTextSearch     = "Search:"
	labelTextTerms      = "Search terms:"
	labelTextMassSearch = "Mass Search:"

	checkTextToggleMassSearch = "Toggle mass search"

	buttonTextAddType = "Add type"
	buttonTextSearch  = "SEARCH"
	buttonTextAddTerm = "Add another\nsearch term"
)

var (
	// TODO: consider enum or mapping for all these
	searchConfTextScope = []string{
		"All Texts",
		"Select Corpus(s)",
		"Select Text(s)",
	}

	searchConfTextSettings = []string{
		"Wrap line text during search",
		"Print text meta information",
		"Use POS instead of EPOS",
		"Use FORM instead of FRAG",
	}

	massSearchConfTextSettings = []string{
		"Use proper noun grouping",
		"Limit search scope",
		"Manual term selection",
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
	labelSearch *widget.Label

	// Horizontal configuration container
	searchConfiguration       *fyne.Container
	searchConfScope           *widget.RadioGroup
	searchConfCorpusSelection *widget.CheckGroup
	searchConfSettings        *widget.CheckGroup

	// Search terms label
	labelSearchTerms *widget.Label
	// Toggle singular vs mass search
	toggleMassSearch *widget.Check
	labelMassSearch  *widget.Label

	// Horizontal search terms container
	terms *fyne.Container

	// Mass search options
	massSearchConfiguration *fyne.Container
	massSearchConfSettings  *widget.CheckGroup

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
	scPane.setupMassSearch()
	scPane.setupSearchButton()

	// Setup labels
	scPane.labelSearch = widget.NewLabel(labelTextSearch)
	scPane.labelSearchTerms = widget.NewLabel(labelTextTerms)
	scPane.labelMassSearch = widget.NewLabel(labelTextMassSearch)

	// Create the finalised search configuration pane
	scPane.mainContainer = container.NewVBox(
		scPane.labelSearch, // Settings label
		widget.NewSeparator(),
		scPane.searchConfiguration, // Search configuration
		widget.NewSeparator(),
		container.NewBorder(
			nil,
			nil,
			scPane.labelSearchTerms, // Terms label
			scPane.toggleMassSearch,
		),
		scPane.terms,        // Search terms
		scPane.termsAdd,     // Add additional search terms
		scPane.searchButton, // Button to actually run the search
	)

	// Save reference to parent to start search
	scPane.gui = gui

	return scPane
}

func (scp *scpStruct) setupConfiguration() {
	// Set up corpus selection so we can enable / disable it depending on selected scope
	scp.searchConfCorpusSelection = widget.NewCheckGroup([]string{}, func(changed []string) {})
	for _, corpusText := range util.ListCorpusAndTexts() { // TODO: expand to select texts in accordion fashion
		scp.searchConfCorpusSelection.Append(corpusText.CorpusName)
	}

	// Radio selection of search scope
	scp.searchConfScope = widget.NewRadioGroup(
		searchConfTextScope,
		func(s string) {
			switch s {
			case searchConfTextScope[0]: // TODO: manually checking position 0 is messy here
				scp.searchConfCorpusSelection.Disable()
			default:
				scp.searchConfCorpusSelection.Enable()
			}
		},
	)
	// Default scope selection of select corpus
	scp.searchConfScope.SetSelected(searchConfTextScope[1]) // TODO: manually checking position 1 is messy here
	// Ensure scope is always set
	scp.searchConfScope.Required = true

	// Search settings
	scp.searchConfSettings = widget.NewCheckGroup(
		searchConfTextSettings,
		func(changed []string) {}, // TODO: is this needed further down the line?
	)

	// Add the search configuration panes in correct order
	scp.searchConfiguration = container.NewHBox(
		scp.searchConfScope,
		widget.NewSeparator(),
		container.NewVScroll(scp.searchConfCorpusSelection),
		widget.NewSeparator(),
		scp.searchConfSettings,
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

func (scp *scpStruct) setupMassSearch() {
	scp.toggleMassSearch = widget.NewCheck(
		checkTextToggleMassSearch,
		func(b bool) {
			// TODO: change from term selection container to mass search container
		},
	)

	scp.massSearchConfSettings = widget.NewCheckGroup(
		massSearchConfTextSettings,
		func(changed []string) {},
	)

	scp.massSearchConfiguration = container.NewHBox(
		scp.massSearchConfSettings,
	)
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
		switch scp.searchConfScope.Selected { // TODO: have this be handled by changing search settings when clicked instead of getting it here
		case searchConfTextScope[0]: // Search all texts
			searchSettings.SearchAllTexts = true
		case searchConfTextScope[1]: // Search selected corpus'
			searchSettings.SelectedCorpus = scp.searchConfCorpusSelection.Selected
		case searchConfTextScope[2]: // Search selected texts // TODO
			searchSettings.SelectedCorpus = scp.searchConfCorpusSelection.Selected
		default: // Should never happen
			panic("cannot determine configuration search scope")
		}
		// TODO: same as above, have it automatically set the settings whenever clicked
		for _, ss := range scp.searchConfSettings.Selected {
			if ss == searchConfTextSettings[0] {
				searchSettings.SearchWrappedText = true
				continue
			}
			if ss == searchConfTextSettings[1] {
				searchSettings.PrintTextMetaInformation = true
				continue
			}
			if ss == searchConfTextSettings[2] {
				searchSettings.UsePOS = true
				continue
			}
			if ss == searchConfTextSettings[3] {
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
