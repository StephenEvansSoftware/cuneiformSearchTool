package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	tabItemNewSearchLabel = "New Search"
)

type searchView interface{}

type svStruct struct {
	sidebar                   *fyne.Container // TODO TODO TODO: make this into a widget.Toolbar
	sidebarButtonAddNewSearch *widget.Button
	sidebarButtonRemoveSearch *widget.Button
	sidebarButtonSaveSearches *widget.Button
	sidebarButtonExportSearch *widget.Button

	searchPane *container.DocTabs

	gui guiInternal
}

func newSearchView(gui guiInternal) fyne.CanvasObject {
	sv := &svStruct{gui: gui}
	sv.setupSidebar()
	sv.setupSearchPane()
	return container.NewBorder(
		nil,
		nil,
		container.NewHBox(
			sv.sidebar,
			widget.NewSeparator(),
		),
		nil,
		sv.searchPane,
	)
}

func (sv *svStruct) setupSidebar() {
	sv.sidebarButtonAddNewSearch = widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {})
	sv.sidebarButtonRemoveSearch = widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {})
	sv.sidebarButtonSaveSearches = widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {})
	sv.sidebarButtonExportSearch = widget.NewButtonWithIcon("", theme.DownloadIcon(), func() {})

	sv.sidebarButtonAddNewSearch.OnTapped = func() {
		sv.searchPane.Append(sv.createNewSearchTab())
	}

	sv.sidebar = container.NewVBox(
		sv.sidebarButtonAddNewSearch,
		widget.NewSeparator(),
		sv.sidebarButtonRemoveSearch,
		widget.NewSeparator(),
		sv.sidebarButtonSaveSearches,
		widget.NewSeparator(),
		sv.sidebarButtonExportSearch,
	)
}

func (sv *svStruct) setupSearchPane() {
	sv.searchPane = container.NewDocTabs(
		sv.createNewSearchTab(),
	)

	sv.searchPane.CreateTab = func() *container.TabItem {
		return sv.createNewSearchTab()
	}

	sv.searchPane.SetTabLocation(container.TabLocationLeading)
}

func (sv *svStruct) createNewSearchTab() *container.TabItem {
	sr := newSearchResultsPane(sv.gui)
	scp := newSearchConfPane(sr)

	searchTab := container.NewTabItem(
		tabItemNewSearchLabel,
		container.NewBorder(
			nil, nil,
			container.NewHBox(
				scp.getPane(),
				widget.NewSeparator(),
			), nil,
			sr.getPane(),
		),
	)
	sr.setParentTab(searchTab)

	return searchTab
}
