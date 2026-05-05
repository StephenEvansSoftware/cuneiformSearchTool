package gui

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"cuneiformSearchTool/util"
)

// This file creates the search configuration pane of the gui

const (
	labelTextSearch     = "Search:"
	labelTextTerms      = "Terms Search:"
	labelTextMassSearch = "Mass Search:"

	checkTextToggleMassSearch = "Toggle mass search"

	buttonTextAddType     = "Add type"
	buttonTextAddTerm     = "Add another\nsearch term"
	buttonTextSearchTerms = "RUN TERMS SEARCH"
	buttonTextSearchMass  = "RUN MASS SEARCH"
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

	startNewSearch()
	finishSearch()
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

	// Lower configuration labels
	labelSearchTerms *widget.Label
	labelMassSearch  *widget.Label
	// Toggle singular vs mass search
	toggleMassSearch *widget.Check

	// Terms container
	terms         *fyne.Container
	termSelection *fyne.Container
	termAdd       *widget.Button

	// Mass search options
	massSearchConfiguration           *fyne.Container
	massSearchConfTermLengthSelection *widget.Select
	massSearchConfInfo                *widget.RichText
	// Easier to handle a set of individual widget.Check than widget.CheckGroup
	massSearchConfSettingProperNounGrouping  *widget.Check
	massSearchConfSettingLimitSearchScope    *widget.Check
	massSearchConfSettingManualTermSelection *widget.Check

	// Buttons to actually run the searches
	searchButtonTerms *widget.Button
	searchButtonMass  *widget.Button

	// Reference to parent gui object to pass information to
	gui guiInternal
}

func newSearchConfigurationPane(gui guiInternal) searchConfigurationPane {
	scPane := &scpStruct{}

	// Setup main pane components
	scPane.setupConfiguration()
	scPane.setupSearchTerms()
	scPane.setupMassSearch()
	scPane.setupSearchTermsButton()
	scPane.setupSearchMassButton()

	// Setup labels
	scPane.labelSearch = widget.NewLabel(labelTextSearch)
	scPane.labelSearchTerms = widget.NewLabel(labelTextTerms)
	scPane.labelMassSearch = widget.NewLabel(labelTextMassSearch)

	// Hide toggleable mass search for now
	scPane.labelMassSearch.Hide()
	scPane.massSearchConfiguration.Hide()
	scPane.searchButtonMass.Hide()

	// Create the finalised search configuration pane
	scPane.mainContainer = container.NewVBox(
		scPane.labelSearch, // Settings label
		widget.NewSeparator(),
		scPane.searchConfiguration, // Search configuration
		widget.NewSeparator(),
		container.NewBorder(
			nil, nil,
			container.NewVBox(
				scPane.labelSearchTerms,
				scPane.labelMassSearch,
			),
			scPane.toggleMassSearch,
		),
		scPane.terms,
		scPane.massSearchConfiguration,
		scPane.searchButtonTerms,
		scPane.searchButtonMass,
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
	scp.termSelection = container.NewHBox()
	scp.termSelection.Add(scp.createNewSearchTerm())
	scp.termAdd = widget.NewButton(buttonTextAddTerm, func() {
		scp.termSelection.Add(scp.createNewSearchTerm())
	})
	scp.terms = container.NewVBox(
		scp.termSelection,
		scp.termAdd,
	)
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
			if b {
				scp.labelSearchTerms.Hide()
				scp.labelMassSearch.Show()
				scp.terms.Hide()
				scp.massSearchConfiguration.Show()
				scp.searchButtonTerms.Hide()
				scp.searchButtonMass.Show()
			} else {
				scp.labelMassSearch.Hide()
				scp.labelSearchTerms.Show()
				scp.massSearchConfiguration.Hide()
				scp.terms.Show()
				scp.searchButtonMass.Hide()
				scp.searchButtonTerms.Show()
			}
		},
	)

	scp.massSearchConfSettingProperNounGrouping = widget.NewCheck(massSearchConfTextSettings[0], func(b bool) {})
	scp.massSearchConfSettingLimitSearchScope = widget.NewCheck(massSearchConfTextSettings[1], func(b bool) {})
	scp.massSearchConfSettingManualTermSelection = widget.NewCheck(massSearchConfTextSettings[2], func(b bool) {})

	scp.massSearchConfInfo = widget.NewRichTextWithText("Total searches\nlisted here\nonce terms\nselected")
	scp.massSearchConfTermLengthSelection = widget.NewSelect(
		[]string{"2", "3", "4", "5"},
		func(s string) {
			massSearchTermsLength, _ := strconv.Atoi(s)
			scp.massSearchConfInfo.ParseMarkdown(fmt.Sprintf("Term length = %d", massSearchTermsLength))
			scp.massSearchConfInfo.AppendMarkdown("Total searches:")
			scp.massSearchConfInfo.AppendMarkdown(fmt.Sprintf("*%d*", util.GetTotalMassSearches(massSearchTermsLength, scp.massSearchConfSettingProperNounGrouping.Checked)))
		},
	)
	scp.massSearchConfTermLengthSelection.PlaceHolder = "Term Length"

	scp.massSearchConfiguration = container.NewHBox(
		container.NewVBox(
			scp.massSearchConfSettingProperNounGrouping,
			scp.massSearchConfSettingLimitSearchScope,
			scp.massSearchConfSettingManualTermSelection,
		),
		scp.massSearchConfTermLengthSelection,
		//		container.NewScroll(scp.massSearchConfInfo),
		scp.massSearchConfInfo,
	)
}

func (scp *scpStruct) setupSearchTermsButton() {
	scp.searchButtonTerms = widget.NewButton(buttonTextSearchTerms, func() {
		// Run a timer to see how long each search takes
		timeStart := time.Now()

		// Collect search terms and settings
		collectedSearchTerms, err := scp.getSearchTerms()
		if err != nil {
			panic(err) // TODO: handle better
		}
		searchSettings, err := scp.getSearchSettings()
		if err != nil {
			panic(err) // TODO: handle better
		}

		// TODO: this isn't quite right, but it'll do for now
		scp.startNewSearch()
		// Inform gui that a new search is happening and let it handle the rest
		fyne.Do(func() {
			scp.gui.startNewTermsSearch(collectedSearchTerms, searchSettings, timeStart)
		})
	})
}

func (scp *scpStruct) setupSearchMassButton() {
	scp.searchButtonMass = widget.NewButton(buttonTextSearchMass, func() {
		// Run a timer to see how long the mass search takes (a long time)
		timeStart := time.Now()

		// Collect search settings and mass search settings
		searchSettings, err := scp.getSearchSettings()
		if err != nil {
			panic(err) // TODO: handle better
		}
		massSearchSettings, err := scp.getMassSearchSettings()
		if err != nil {
			panic(err) // TODO: handle better
		}

		scp.startNewSearch()
		fyne.Do(func() {
			scp.gui.startNewMassSearch(searchSettings, massSearchSettings, timeStart)
		})
	})
}

func (scp *scpStruct) getSearchTerms() ([][]string, error) {
	// Retrieve all search terms
	collectedSearchTerms := [][]string{}
	for sti, st := range scp.termSelection.Objects { // Overarching search terms container
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
							return nil, fmt.Errorf("term contains non-expected type: %v", toot)
						}
					}
				}
			case *widget.Button: // "Add type" button, ignore
			default: // Unknown type, should never happen
				// TODO: handle better?
				return nil, fmt.Errorf("terms contain non-expected type: %v", sto)
			}
		}
	}
	return collectedSearchTerms, nil
}

// TODO: take from above method and put in here
func (scp *scpStruct) getSearchSettings() (*util.SearchSettings, error) {
	searchSettings := &util.SearchSettings{}
	switch scp.searchConfScope.Selected { // TODO: have this be handled by changing search settings when clicked instead of getting it here
	case searchConfTextScope[0]: // Search all texts
		searchSettings.SearchAllTexts = true
	case searchConfTextScope[1]: // Search selected corpus'
		searchSettings.SelectedCorpus = scp.searchConfCorpusSelection.Selected
	case searchConfTextScope[2]: // Search selected texts // TODO
		searchSettings.SelectedCorpus = scp.searchConfCorpusSelection.Selected
	default: // Should never happen
		return nil, fmt.Errorf("cannot determine configuration search scope")

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
	return searchSettings, nil
}

func (scp *scpStruct) getMassSearchSettings() (*util.MassSearchSettings, error) {
	massSearchTermsLength, _ := strconv.Atoi(scp.massSearchConfTermLengthSelection.Selected)
	massSearchSettings := &util.MassSearchSettings{
		TermsLength:           massSearchTermsLength,
		UseProperNounGrouping: scp.massSearchConfSettingProperNounGrouping.Checked,
		LimitSearchScope:      scp.massSearchConfSettingLimitSearchScope.Checked,
		ManualTermSelection:   scp.massSearchConfSettingManualTermSelection.Checked,
	}

	return massSearchSettings, nil
}

func (scp *scpStruct) getPane() fyne.CanvasObject {
	return scp.mainContainer
}

// Disables widgets so nothing can be changed / clicked while search is running
func (scp *scpStruct) startNewSearch() {
	scp.searchConfScope.Disable()
	scp.searchConfCorpusSelection.Disable()
	scp.searchConfSettings.Disable()
	scp.toggleMassSearch.Disable()
	scp.termAdd.Disable()
	scp.massSearchConfSettingProperNounGrouping.Disable()
	scp.massSearchConfSettingLimitSearchScope.Disable()
	scp.massSearchConfSettingManualTermSelection.Disable()
	scp.searchButtonTerms.Disable()
	scp.searchButtonMass.Disable()
}

// Re-enables all widgets
func (scp *scpStruct) finishSearch() {
	scp.searchConfScope.Enable()
	scp.searchConfCorpusSelection.Enable()
	scp.searchConfSettings.Enable()
	scp.toggleMassSearch.Enable()
	scp.termAdd.Enable()
	scp.massSearchConfSettingProperNounGrouping.Enable()
	scp.massSearchConfSettingLimitSearchScope.Enable()
	scp.massSearchConfSettingManualTermSelection.Enable()
	scp.searchButtonTerms.Enable()
	scp.searchButtonMass.Enable()
}
