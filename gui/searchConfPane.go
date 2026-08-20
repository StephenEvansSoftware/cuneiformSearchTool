package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"cuneiformSearchTool/util"
)

const (
	labelTextSearch               = "# **Search Settings:**"
	labelTextSearchSelectCorpuses = "Select Corpuses:"
	labelTextSearchSelectTexts    = "Select Texts:"
	labelTextSearchSettings       = "Toggleable Settings:"
	labelTextSearchTerms          = "Search Terms:"

	selectCorpusAll    = "All Corpuses"
	selectCorpusManual = "Select Corpuses"

	selectTextsAll    = "All Texts"
	selectTextsManual = "Select Texts"

	settingWrapLines     = "Wrap lines"
	settingIgnoreOverlap = "Ignore overlapping terms"
	settingUsePOS        = "Use POS instead of EPOS"
	settingShowFORM      = "Show FORM in result instead of FRAG"

	termAddNew = "Add additional term"

	buttonRunSearch = "Run Search"
)

type searchConfPane interface {
	paneGetter
}

type scpStruct struct {
	mainContainer *fyne.Container

	searchScope               *fyne.Container
	searchScopeCorpuses       *widget.RadioGroup
	searchScopeCorpusesManual *widget.CheckGroup
	searchScopeTexts          *widget.RadioGroup
	searchScopeTextsManual    *widget.CheckGroup

	searchSettings       *fyne.Container
	searchSettingsChecks *widget.CheckGroup

	searchTerms    *fyne.Container
	searchTermsAll *fyne.Container

	searchButton *widget.Button

	srp searchResultsPane
}

func newSearchConfPane(srp searchResultsPane) searchConfPane {
	scp := &scpStruct{srp: srp}

	scp.setupSearchScope()
	scp.setupSearchSettings()
	scp.setupSearchTerms()
	scp.setupSearchButton()

	scp.mainContainer = container.NewBorder(
		nil, scp.searchButton,
		nil, nil,
		container.NewVScroll(
			container.NewVBox(
				widget.NewRichTextFromMarkdown(labelTextSearch),
				widget.NewSeparator(),
				scp.searchScope,
				scp.searchSettings,
				scp.searchTerms,
			),
		),
	)

	return scp
}

func (scp *scpStruct) setupSearchScope() {
	labelSelectCorpuses := widget.NewLabelWithStyle(labelTextSearchSelectCorpuses, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	labelSelectTexts := widget.NewLabelWithStyle(labelTextSearchSelectTexts, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	scp.searchScopeCorpusesManual = widget.NewCheckGroup([]string{}, func(c []string) {})
	scp.searchScopeTextsManual = widget.NewCheckGroup([]string{}, func(c []string) {})
	for _, ct := range util.ListCorpusAndTexts() {
		scp.searchScopeCorpusesManual.Append(ct.CorpusName)
		scp.searchScopeTextsManual.Append(ct.CorpusName)
	}

	scp.searchScopeCorpuses = widget.NewRadioGroup(
		[]string{selectCorpusAll, selectCorpusManual},
		func(s string) {
			switch s {
			case selectCorpusAll:
				scp.searchScopeCorpusesManual.Disable()
			case selectCorpusManual:
				scp.searchScopeCorpusesManual.Enable()
			}
		},
	)
	scp.searchScopeCorpuses.SetSelected(selectCorpusAll)
	scp.searchScopeCorpuses.Required = true

	scp.searchScopeTexts = widget.NewRadioGroup(
		[]string{selectTextsAll, selectTextsManual},
		func(s string) {
			switch s {
			case selectTextsAll:
				scp.searchScopeTextsManual.Disable()
			case selectTextsManual:
				scp.searchScopeTextsManual.Enable()
			}
		},
	)
	scp.searchScopeTexts.SetSelected(selectTextsAll)
	scp.searchScopeTexts.Required = true

	scp.searchScope = container.NewVBox(
		labelSelectCorpuses,
		container.NewBorder(
			nil, nil,
			scp.searchScopeCorpuses,
			container.NewVScroll(scp.searchScopeCorpusesManual), // TODO: add padding so three can be seen
		),
		labelSelectTexts,
		container.NewBorder(
			nil, nil,
			scp.searchScopeTexts,
			container.NewVScroll(scp.searchScopeTextsManual), // TODO: add padding so three can be seen
		),
	)
}

func (scp *scpStruct) setupSearchSettings() {
	settingsLabel := widget.NewLabelWithStyle(labelTextSearchSettings, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	scp.searchSettingsChecks = widget.NewCheckGroup(
		[]string{settingWrapLines, settingIgnoreOverlap, settingUsePOS, settingShowFORM},
		func(c []string) {},
	)

	scp.searchSettings = container.NewVBox(
		settingsLabel,
		scp.searchSettingsChecks,
	)
}

func (scp *scpStruct) setupSearchTerms() {
	termsLabel := widget.NewLabelWithStyle(labelTextSearchTerms, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	scp.searchTermsAll = container.NewVBox()
	scp.searchTermsAll.Add(scp.createNewSearchTerm())

	scp.searchTerms = container.NewVBox(
		termsLabel,
		scp.searchTermsAll,
		widget.NewSeparator(),
		widget.NewButton(termAddNew, func() {
			scp.searchTermsAll.Add(scp.createNewSearchTerm())
		}),
	)

}

func (scp *scpStruct) createNewSearchTerm() *fyne.Container {
	searchTerm := container.NewVBox(widget.NewSeparator())
	searchTypes := container.NewVBox()

	createType := func() fyne.CanvasObject {
		removeType := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {})
		newType := container.NewBorder(
			nil, nil, nil,
			removeType,
			widget.NewSelect(append([]string{util.AllProperNounsType}, util.ValidTypes...), func(s string) {}),
		)
		removeType.OnTapped = func() {
			if len(searchTypes.Objects) == 1 && len(scp.searchTermsAll.Objects) == 1 {
				return
			}
			searchTypes.Remove(newType)
			if len(searchTypes.Objects) == 0 {
				scp.searchTermsAll.Remove(searchTerm)
			}
		}
		return newType
	}

	searchTypes.Add(createType())

	addType := widget.NewButtonWithIcon("Add new type", theme.ContentAddIcon(), func() {
		searchTypes.Add(createType())
	})

	searchTerm.Add(searchTypes)
	searchTerm.Add(addType)
	searchTerm.Add(widget.NewSeparator())

	return searchTerm
}

func (scp *scpStruct) setupSearchButton() {
	scp.searchButton = widget.NewButton(buttonRunSearch, func() {
		searchTimeStart := time.Now()
		termsToSearch := scp.getSearchTerms()
		settingsForSearch := scp.getSearchSettings()
		settingsForSearch.TermsLength = len(termsToSearch)

		scp.disable()

		err := fmt.Errorf("")
		if err = util.ValidateTypes(termsToSearch); err != nil {
			scp.srp.setError(err.Error())
			scp.enable()
			return
		}

		scp.srp.searchStarted(settingsForSearch, termsToSearch)
		searchResults := []*util.CorpusSearchResult{}
		totalInstances := -1
		fyne.Do(func() {
			searchResults, totalInstances, err = util.SearchFor(
				termsToSearch,
				settingsForSearch,
			)

			if err != nil {
				scp.srp.setError(err.Error())
				scp.enable()
				return
			}

			searchTimeDuration := time.Since(searchTimeStart)

			scp.srp.searchFinished(searchResults, totalInstances, searchTimeStart, searchTimeDuration)
			scp.mainContainer.Hide()
		})
	})
}

func (scp *scpStruct) getSearchTerms() [][]string {
	collectedSearchTerms := [][]string{}
	for sti, st := range scp.searchTermsAll.Objects {
		collectedSearchTerms = append(collectedSearchTerms, []string{})
		for _, sto := range st.(*fyne.Container).Objects {
			switch t := sto.(type) {
			case *fyne.Container:
				for _, to := range t.Objects {
					for _, too := range to.(*fyne.Container).Objects {
						switch toot := too.(type) {
						case *widget.Select:
							if selected := toot.Selected; selected != "" {
								collectedSearchTerms[sti] = append(collectedSearchTerms[sti], selected)
							}
						}
					}
				}
			}
		}
	}
	return collectedSearchTerms
}

func (scp *scpStruct) getSearchSettings() *util.SearchSettings {
	searchSettings := &util.SearchSettings{}
	switch scp.searchScopeCorpuses.Selected {
	case selectCorpusAll:
		searchSettings.SearchAllCorpuses = true
	case selectCorpusManual:
		searchSettings.SelectedCorpus = scp.searchScopeCorpusesManual.Selected
	}

	switch scp.searchScopeTexts.Selected {
	case selectTextsAll:
		searchSettings.SearchAllTexts = true
	case selectTextsManual:
		searchSettings.SelectedTexts = scp.searchScopeTextsManual.Selected
	}

	for _, ss := range scp.searchSettingsChecks.Selected {
		switch ss {
		case settingWrapLines:
			searchSettings.SearchWrappedText = true
		case settingIgnoreOverlap:
			searchSettings.IgnoreOverlap = true
		case settingUsePOS:
			searchSettings.UsePOS = true
		case settingShowFORM:
			searchSettings.UseFORM = true
		}
	}

	return searchSettings
}

func (scp *scpStruct) disable() {
	scp.searchScopeCorpuses.Disable()
	scp.searchScopeCorpusesManual.Disable()
	scp.searchScopeTexts.Disable()
	scp.searchScopeTextsManual.Disable()
	scp.searchSettingsChecks.Disable()
	// TODO: horrible, make better
	for _, obj := range scp.searchTerms.Objects {
		switch ot := obj.(type) {
		case *widget.Button:
			ot.Disable()
		case *fyne.Container:
			for _, otc := range ot.Objects {
				switch otct := otc.(type) {
				case *fyne.Container:
					for _, otctc := range otct.Objects {
						switch otctct := otctc.(type) {
						case *widget.Button:
							otctct.Disable()
						case *fyne.Container:
							for _, otctctc := range otctct.Objects {
								switch otctctct := otctctc.(type) {
								case *fyne.Container:
									for _, otctctctc := range otctctct.Objects {
										switch otctctctct := otctctctc.(type) {
										case *widget.Button:
											otctctctct.Disable()
										case *widget.Select:
											otctctctct.Disable()
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func (scp *scpStruct) enable() {
	scp.searchScopeCorpuses.Enable()
	scp.searchScopeCorpusesManual.Enable()
	scp.searchScopeTexts.Enable()
	scp.searchScopeTextsManual.Enable()
	scp.searchSettingsChecks.Enable()
	// TODO: horrible, make better
	for _, obj := range scp.searchTerms.Objects {
		switch ot := obj.(type) {
		case *widget.Button:
			ot.Enable()
		case *fyne.Container:
			for _, otc := range ot.Objects {
				switch otct := otc.(type) {
				case *fyne.Container:
					for _, otctc := range otct.Objects {
						switch otctct := otctc.(type) {
						case *widget.Button:
							otctct.Enable()
						case *fyne.Container:
							for _, otctctc := range otctct.Objects {
								switch otctctct := otctctc.(type) {
								case *fyne.Container:
									for _, otctctctc := range otctctct.Objects {
										switch otctctctct := otctctctc.(type) {
										case *widget.Button:
											otctctctct.Enable()
										case *widget.Select:
											otctctctct.Enable()
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func (scp *scpStruct) getPane() fyne.CanvasObject {
	return scp.mainContainer
}
